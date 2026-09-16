// Copyright The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func TestRefreshDevicesRespectsConcurrency(t *testing.T) {
	devices := []Device{
		{Name: "one"},
		{Name: "two"},
		{Name: "three"},
		{Name: "four"},
	}
	logger := slog.New(slog.DiscardHandler)

	var mutex sync.Mutex
	running := 0
	maxRunning := 0
	refreshDevices(logger, devices, 1, func(_ *slog.Logger, _ Device) {
		mutex.Lock()
		running++
		if running > maxRunning {
			maxRunning = running
		}
		mutex.Unlock()

		time.Sleep(10 * time.Millisecond)

		mutex.Lock()
		running--
		mutex.Unlock()
	})

	if maxRunning != 1 {
		t.Fatalf("maximum concurrent reads = %d, want 1", maxRunning)
	}
}

func TestShouldRefreshAfterLongRefresh(t *testing.T) {
	// reset global interval
	originalInterval := *smartctlInterval
	*smartctlInterval = time.Minute
	t.Cleanup(func() { *smartctlInterval = originalInterval })

	completed := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	collector := SMARTctlManagerCollector{
		lastRefreshCompleted: completed,
		lastRefreshDuration:  time.Minute,
	}

	if collector.shouldRefresh(completed.Add(59 * time.Second)) {
		t.Fatal("started a new refresh before the next interval")
	}
	if !collector.shouldRefresh(completed.Add(time.Minute)) {
		t.Fatal("did not start a refresh at the next interval")
	}
}

func TestCollectDoesNotStartRefreshWhileAnotherIsRunning(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	firstRefreshStarted := make(chan struct{})
	releaseFirstRefresh := make(chan struct{})
	var refreshes atomic.Int32
	collector := SMARTctlManagerCollector{
		logger: logger,
		refresh: func(_ *slog.Logger, _ []Device) {
			if refreshes.Add(1) == 1 {
				close(firstRefreshStarted)
				<-releaseFirstRefresh
			}
		},
	}

	firstDone := make(chan struct{})
	go func() {
		collector.Collect(make(chan prometheus.Metric, 1))
		close(firstDone)
	}()
	<-firstRefreshStarted

	secondDone := make(chan struct{})
	go func() {
		collector.Collect(make(chan prometheus.Metric, 1))
		close(secondDone)
	}()

	select {
	case <-secondDone:
		t.Fatal("second collection completed while the first refresh was running")
	case <-time.After(20 * time.Millisecond):
	}
	if got := refreshes.Load(); got != 1 {
		t.Fatalf("refreshes started = %d, want 1", got)
	}

	close(releaseFirstRefresh)
	<-firstDone
	<-secondDone
}
