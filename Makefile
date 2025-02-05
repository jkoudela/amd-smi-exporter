.PHONY: build build-linux clean

# Binary name
BINARY_NAME=amd-smi-exporter

# Build the binary for the local architecture (macOS)
build:
	go build -o $(BINARY_NAME)

# Build for Linux (amd64)
build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-linux-amd64
