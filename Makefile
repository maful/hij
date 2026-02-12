.PHONY: build run clean test

BINARY_NAME=hij
BUILD_DIR=build

build:
	@mkdir -p $(BUILD_DIR)
	@echo "Building $(BINARY_NAME)..."
	@go build -ldflags "-s -w -X main.version=$$(git describe --tags --always --dirty) -X main.commit=$$(git rev-parse --short HEAD) -X main.date=$$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o $(BUILD_DIR)/$(BINARY_NAME) .

run:
	go run .

clean:
	rm -f $(BUILD_DIR)/$(BINARY_NAME)

test:
	go test -v ./...
