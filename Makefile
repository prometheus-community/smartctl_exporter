GOPATH=$(shell pwd)/vendor:$(shell pwd)
GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=$(shell basename "$(PWD)")

build: get
	@echo "Building $(GOFILES) to ./bin"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -v -o bin/$(GONAME) $(GOFILES)

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
