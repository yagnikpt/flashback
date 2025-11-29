BINARY_NAME=flashback
SRC_DIR=.
GOFLAGS=-ldflags="-s -w"

.PHONY: all
all: build

.PHONY: build
build:
	go build $(GOFLAGS) -o $(BINARY_NAME) $(SRC_DIR)

.PHONY: run
run:
	go run $(SRC_DIR)

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: install
install: build
	mv $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
