# Makefile for odcread (Go version)

BINARY_NAME=odcread
SRC_DIR=src
BUILD_DIR=bin
TEST_DIR=_tests

.PHONY: all build clean test fmt lint run help

all: build

help:
	@echo "Usage:"
	@echo "  make build    - Build the odcread binary"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make test     - Run integration tests against .odc files in $(TEST_DIR)"
	@echo "  make fmt      - Format Go source code"
	@echo "  make lint     - Run go vet on source code"
	@echo "  make run FILE=path/to/file.odc - Run odcread on a specific file"

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@cd $(SRC_DIR) && go build -o ../$(BUILD_DIR)/$(BINARY_NAME) ./cmd/odcread

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

test: build
	@echo "Running integration tests..."
	@for f in $(TEST_DIR)/mini*.odc; do \
		echo "=== $$(basename $$f) ==="; \
		./$(BUILD_DIR)/$(BINARY_NAME) "$$f" 2>&1 | grep -v "Read store" || true; \
	done

fmt:
	@echo "Formatting code..."
	@cd $(SRC_DIR) && go fmt ./...

lint:
	@echo "Linting code..."
	@cd $(SRC_DIR) && go vet ./...

run: build
	@if [ -z "$(FILE)" ]; then \
		echo "Please provide a file: make run FILE=path/to/file.odc"; \
		exit 1; \
	fi
	./$(BUILD_DIR)/$(BINARY_NAME) "$(FILE)"
