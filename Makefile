# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
SERVER_BINARY=server
CLIENT_BINARY=client

# Build directories
BUILD_DIR=build
SERVER_DIR=cmd/server
CLIENT_DIR=cmd/client

.PHONY: all build clean test server client

all: clean build

build: server client

server:
	@echo "Building server..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(SERVER_BINARY) ./$(SERVER_DIR)

client:
	@echo "Building client..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(CLIENT_BINARY) ./$(CLIENT_DIR)

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

test:
	$(GOTEST) -v ./...

deps:
	$(GOGET) -v ./...
