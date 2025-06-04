.PHONY: lint test build clean

# Build the application
build:
	go build -o golang-kata-1 .

# Run all tests
test:
	go test -count 1 -p 1 -failfast -v ./...

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -f golang-kata-1

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
bench:
	go test -bench=. -benchmem ./...

# Format code
fmt:
	go fmt ./...

# Run all checks (lint, test, build)
check: lint test build
