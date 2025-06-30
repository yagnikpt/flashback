# Variables
BINARY_NAME=flashback
SRC_DIR=./cmd/flashback
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
tidy::
	go mod tidy -v
