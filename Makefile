GOPATH=$(shell pwd)/vendor:$(shell pwd)
GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=$(shell basename "$(PWD)")
PID=/tmp/go-$(GONAME).pid

build: get
	@echo "Building $(GOFILES) to ./bin"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o bin/$(GONAME) $(GOFILES)

get:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get .

install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

run: build
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go run $(GOFILES) --config=$(shell pwd)/smartctl_exporter.yaml --debug --verbose

run-sudo: build
	sudo bin/$(GONAME) --config=$(shell pwd)/smartctl_exporter.yaml --debug --verbose

watch:
	@$(MAKE) restart &
	@fswatch -o . -e 'bin/.*' | xargs -n1 -I{}  make restart

restart: clear stop clean build start

start: build
	@echo "Starting bin/$(GONAME)"
	@./bin/$(GONAME) & echo $$! > $(PID)

stop:
	@echo "Stopping bin/$(GONAME) if it's running"
	@-kill `[[ -f $(PID) ]] && cat $(PID)` 2>/dev/null || true

clear:
	@clear

clean:
	@echo "Cleaning"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
