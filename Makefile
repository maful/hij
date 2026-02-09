.PHONY: build run clean

BINARY_NAME=hij
BUILD_DIR=build

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

run:
	go run .

clean:
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
