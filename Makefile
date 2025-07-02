# Variables
BINARY_NAME=flashback
SRC_DIR=./cmd/flashback
GOFLAGS=-ldflags="-s -w"
GOLANG_CROSS_VERSION  ?= v1.24.4

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

.PHONY: release-dry-run
release-dry-run:
	@podman run \
		--rm \
		-e CGO_ENABLED=1 \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--clean --skip=validate --skip=publish
