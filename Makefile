# Variables
BINARY_NAME=whrelay
SERVER_BINARY=whrelay_server
SERVER_DIR=./server
CLI_DIR=./cmd/whrelay
GOFLAGS=-ldflags="-s -w"
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

.DEFAULT_GOAL := all

.PHONY: all
all: build-cli build-server

.PHONY: build-cli
build-cli:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GOFLAGS) -o $(BINARY_NAME) $(CLI_DIR)

.PHONY: build-server
build-server:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GOFLAGS) -o $(SERVER_DIR)/$(SERVER_BINARY) $(SERVER_DIR)

.PHONY: cli
cli:
	go run $(CLI_DIR)

.PHONY: server
server:
	go run $(SERVER_DIR)

.PHONY: download
download:
	go mod download

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -f $(BINARY_NAME) $(SERVER_DIR)/$(SERVER_BINARY)