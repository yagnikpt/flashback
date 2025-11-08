# Variables
BINARY_NAME=flashback
DAEMON_BINARY_NAME=flashbackd
SRC_DIR=./cmd/flashback
DAEMON_DIR=./cmd/daemon
GOFLAGS=-ldflags="-s -w"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	go build $(GOFLAGS) -o $(BINARY_NAME) $(SRC_DIR)

# Build the daemon
# .PHONY: build-daemon
# build-daemon:
# 	go build $(GOFLAGS) -o $(DAEMON_BINARY_NAME) $(DAEMON_DIR)

# Run the application
.PHONY: run
run: build
	./$(BINARY_NAME)

# Run tidy
.PHONY: tidy
tidy:
	go mod tidy -v
