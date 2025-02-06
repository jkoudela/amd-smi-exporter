.PHONY: build clean

# Binary name
BINARY_NAME=amd-smi-exporter

# Build for Linux (amd64)
build:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
