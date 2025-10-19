# Variables
BINARY_NAME=whrelay
SRC_DIR=./server
GOFLAGS=-ldflags="-s -w"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	go build $(GOFLAGS) -o $(BINARY_NAME) $(SRC_DIR)

# Run the application
.PHONY: run
run:
	go run $(SRC_DIR)

# Run tidy
.PHONY: tidy
tidy:
	go mod tidy -v