.PHONY: build clean update-dashboard update-metrics tidy test

# Binary name
BINARY_NAME=amd-smi-exporter

# Build for Linux (amd64)
build:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)

# Update Grafana dashboard
update-dashboard:
	./scripts/update-dashboard.sh

update-metrics:
	./scripts/update-metrics.sh

# Add dependencies
tidy:
	go mod tidy

# Run tests
test:
	go test ./...
