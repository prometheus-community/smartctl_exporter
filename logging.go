// Copyright 2022 The Prometheus Authors
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
	"fmt"
)

// Logger is logger
type Logger struct {
	verbose bool
	debug   bool
}

func newLogger(verbose bool, debug bool) Logger {
	logger := Logger{verbose: verbose, debug: debug}
	return logger
}

// Msg formatted message
func (i Logger) print(prefix string, format string, values ...interface{}) {
	fmt.Printf(fmt.Sprintf("[%s] %s\n", prefix, format), values...)
}

// Info log message
func (i Logger) Info(format string, values ...interface{}) {
	i.print("Info", format, values...)
}

// Warning log message
func (i Logger) Warning(format string, values ...interface{}) {
	i.print("Warning", format, values...)
}

// Error log message
func (i Logger) Error(format string, values ...interface{}) {
	i.print("Error", format, values...)
}

// Panic log message
func (i Logger) Panic(format string, values ...interface{}) {
	i.print("Panic", format, values...)
}

// Verbose log message
func (i Logger) Verbose(format string, values ...interface{}) {
	if i.verbose {
		i.print("Verbose", format, values...)
	}
}

// Debug log message
func (i Logger) Debug(format string, values ...interface{}) {
	if i.debug {
		i.print("Debug", format, values...)
	}
}
