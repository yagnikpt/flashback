# Variables
BINARY_NAME=flashback
DAEMON_BINARY_NAME=flashback-daemon
PACKAGE_NAME          := github.com/yagnik-patel-47/flashback
SRC_DIR=./cmd/flashback
DAEMON_DIR=./cmd/daemon
GOFLAGS=-ldflags="-s -w"
GOLANG_CROSS_VERSION  ?= v1.24.4

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	go build $(GOFLAGS) -o $(BINARY_NAME) $(SRC_DIR)

# Build the daemon
.PHONY: build-daemon
build-daemon:
	go build $(GOFLAGS) -o $(DAEMON_BINARY_NAME) $(DAEMON_DIR)

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
		-e GOCACHE=/tmp/.cache \
		-v /run/user/1000/podman/podman.sock:/run/user/1000/podman/podman.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME):Z \
		-v `pwd`/.cache:/tmp/.cache:Z \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean --skip=validate --skip=publish --snapshot
