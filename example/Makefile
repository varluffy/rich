APPNAME = application
VERSION=$(shell git describe --tags --always)
GOBIN=$(shell go env GOBIN)

.PHONY: init
init:
	go get -u github.com/google/wire/cmd/wire

.PHONY: wire
wire:
	cd cmd/application/ && wire .
.PHONY: run
run:
	cd cmd/application/ && go run . -c ../../config/config.yaml
.PHONY: build
build:
	go build -ldflags "-X main.Version=$(VERSION) -X main.Name=$(APPNAME)" -o $(GOBIN)/ ./...
.PHONY: run
test:
	go test -v ./... -cover