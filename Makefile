# Binary name
BINARY_NAME=wazuh-cli

# Build directory
BUILD_DIR=bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Versioning (git tags or 'dev')
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X 'github.com/ba0f3/wazuh-cli/cmd.Version=$(VERSION)'"

.PHONY: all build clean test install lint tidy fmt vet help

all: clean test build

build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

install:
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) .

lint:
	golangci-lint run ./...

tidy:
	$(GOMOD) tidy

fmt:
	$(GOFMT) ./...

vet:
	$(GOCMD) vet ./...

help:
	@echo "Makefile for wazuh-cli"
	@echo ""
	@echo "Usage:"
	@echo "  make build    - Build the binary"
	@echo "  make test     - Run unit tests"
	@echo "  make vet      - Run go vet"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make install  - Install binary to GOPATH/bin"
	@echo "  make tidy     - Run go mod tidy"
	@echo "  make fmt      - Format Go source code"
	@echo "  make lint     - Run golangci-lint (if installed)"
	@echo ""
