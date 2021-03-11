APPNAME = application
VERSION=$(shell git describe --tags --always)
GOBIN=$(shell go env GOBIN)

.PHONY: init
init:
	go get -u github.com/google/wire/cmd/wire

.PHONY: wire
wire:
	cd cmd/server/ && wire .
.PHONY: run
run:
	cd cmd/server/ && go run . -c ../../config/config.yaml
.PHONY: build
build:
	go build -ldflags "-X main.Version=$(VERSION) -X main.Name=$(APPNAME)" -o $(GOBIN)/ ./...
.PHONY: run
test:
	go test -v ./... -cover
.PHONY: gen
gen:
	protoc -I ./api \
	--openapiv2_out ./api --openapiv2_opt logtostderr=true \
	--openapiv2_opt json_names_for_fields=false \
	--go_out ./api --go_opt=paths=source_relative \
	--go-gin_out ./api --go-gin_opt=paths=source_relative \
	api/app/v1/v1.proto
	protoc-go-inject-tag -input=./api/app/v1/v1.pb.go