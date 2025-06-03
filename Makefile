# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
SERVER_BINARY=server
CLIENT_BINARY=client
WEBSERVER_BINARY=webserver

# Build directories
BUILD_DIR=build
SERVER_DIR=cmd/server
CLIENT_DIR=cmd/client
WEBSERVER_DIR=cmd/webserver

.PHONY: all build clean test server client webserver

all: clean build

build: server client webserver

server:
	@echo "Building server..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(SERVER_BINARY) ./$(SERVER_DIR)

client:
	@echo "Building client..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(CLIENT_BINARY) ./$(CLIENT_DIR)

webserver:
	@echo "Building webserver..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(WEBSERVER_BINARY) ./$(WEBSERVER_DIR)

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

test:
	$(GOTEST) -v ./...

deps:
	$(GOGET) -v ./...
