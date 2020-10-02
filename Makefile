build: get
	@go build -v

get:
	@go get -v

install:
	@go install

run: build
	@go run . --config=$(shell pwd)/smartctl_exporter.yaml --debug --verbose

run-sudo: build
	sudo ./smartctl_exporter --config=$(shell pwd)/smartctl_exporter.yaml --debug --verbose

clear:
	@clear

clean:
	@echo "Cleaning"
	@go clean

example:
	@echo "# Example output" > EXAMPLE.md
	@echo '```' >> EXAMPLE.md
	@curl -s localhost:9633/metrics | grep smartctl >> EXAMPLE.md
	@echo '```' >> EXAMPLE.md
