GOPATH=$(shell pwd)/vendor:$(shell pwd)
GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=$(shell basename "$(PWD)")

build: get
	@echo "Building $(GOFILES) to ./bin"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -v -o bin/$(GONAME) $(GOFILES)

build-static:
	@echo "Building $(GOFILES) to ./bin"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -a -tags netgo -ldflags '-w -extldflags "-static"' -o bin/smartctl_exporter_static *.go

get:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get -v .

install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

run: build
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go run $(GOFILES) --config=$(shell pwd)/smartctl_exporter.yaml --debug --verbose

run-sudo: build
	sudo bin/$(GONAME) --config=$(shell pwd)/smartctl_exporter.yaml --debug --verbose

clear:
	@clear

clean:
	@echo "Cleaning"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

example:
	@echo "# Example output" > EXAMPLE.md
	@echo '```' >> EXAMPLE.md
	@curl -s localhost:9633/metrics | grep smartctl >> EXAMPLE.md
	@echo '```' >> EXAMPLE.md

collect_fake_json:
	-mkdir debug
	-rm -f debug/*json
	sudo ./collect_fake_json.sh
