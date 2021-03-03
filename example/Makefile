APPNAME = application
VERSION=$(shell git describe --tags --always)
GOBIN=$(shell go env GOBIN)

.PHONY: run
run:
	cd cmd/application/ && go run . -c ../../config/config.yaml
.PHONY: build
build:
	go build -ldflags "-X main.Version=$(VERSION) -X main.Name=$(APPNAME)" -o $(GOBIN)/ ./...
.PHONY: run
test:
	go test -v ./... -cover