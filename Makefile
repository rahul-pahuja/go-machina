# Makefile for GoMachina

# Variables
BINARY_NAME=server
BINARY_DIR=bin
MAIN_FILE=cmd/server/main.go
PACKAGES=./machina
COVER_PROFILE=coverage.out
COVER_HTML=coverage.html

# Default target
.DEFAULT_GOAL := help

# Build the server application
build:
	go build -o $(BINARY_DIR)/$(BINARY_NAME) $(MAIN_FILE)

# Build the server application with race detection
build-race:
	go build -race -o $(BINARY_DIR)/$(BINARY_NAME) $(MAIN_FILE)

# Run the server application
run:
	go run $(MAIN_FILE)

# Run the server application with race detection
run-race:
	go run -race $(MAIN_FILE)

# Run tests
test:
	go test ./...

# Run tests without cache
test-nocache:
	go test -count=1 ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with verbose output without cache
test-verbose-nocache:
	go test -v -count=1 ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Run tests with coverage without cache
test-coverage-nocache:
	go test -cover -count=1 ./...

# Run tests with coverage profile
test-coverage-profile:
	go test -coverprofile=$(COVER_PROFILE) ./...

# Run tests with coverage profile without cache
test-coverage-profile-nocache:
	go test -coverprofile=$(COVER_PROFILE) -count=1 ./...

# Run tests with coverage and generate HTML report
test-coverage-html: test-coverage-profile
	go tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)

# Run tests with coverage and generate HTML report without cache
test-coverage-html-nocache:
	go test -coverprofile=$(COVER_PROFILE) -count=1 ./...
	go tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)

# Run tests with race detection
test-race:
	go test -race ./...

# Run tests with race detection without cache
test-race-nocache:
	go test -race -count=1 ./...

# Run benchmarks
benchmark:
	go test -bench=. -benchmem ./...

# Run benchmarks without cache
benchmark-nocache:
	go test -bench=. -benchmem -count=1 ./...

# Run benchmarks with race detection
benchmark-race:
	go test -bench=. -benchmem -race ./...

# Run benchmarks with race detection without cache
benchmark-race-nocache:
	go test -bench=. -benchmem -race -count=1 ./...

# Tidy go.mod
tidy:
	go mod tidy

# Install dependencies
install:
	go mod download

# Clean build artifacts
clean:
	rm -rf $(BINARY_DIR)/
	rm -f $(COVER_PROFILE) $(COVER_HTML)

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Check for style issues using golangci-lint (if installed)
lint:
	golangci-lint run ./...

# Generate mocks
generate:
	go generate ./...

# Run all checks (vet, fmt, test)
check: vet fmt test

# Run all checks without cache
check-nocache: vet fmt test-nocache

# Run all checks with race detection
check-race: vet fmt test-race

# Run all checks with race detection without cache
check-race-nocache: vet fmt test-race-nocache

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Update dependencies
update:
	go get -u ./...

# Show test coverage
coverage: test-coverage-profile
	go tool cover -func=$(COVER_PROFILE)

# Show test coverage without cache
coverage-nocache: test-coverage-profile-nocache
	go tool cover -func=$(COVER_PROFILE)

# Generate documentation
godoc:
	godoc -http=:6060

# Run integration tests
test-integration:
	go test -tags=integration ./...

# Run integration tests without cache
test-integration-nocache:
	go test -tags=integration -count=1 ./...

# Run all tests including integration tests
test-all: test test-integration

# Run all tests including integration tests without cache
test-all-nocache: test-nocache test-integration-nocache

# Help
help:
	@echo "Available commands:"
	@echo "  build                      - Build the server application"
	@echo "  build-race                 - Build the server application with race detection"
	@echo "  run                        - Run the server application"
	@echo "  run-race                   - Run the server application with race detection"
	@echo "  test                       - Run tests"
	@echo "  test-nocache               - Run tests without cache"
	@echo "  test-verbose               - Run tests with verbose output"
	@echo "  test-verbose-nocache       - Run tests with verbose output without cache"
	@echo "  test-coverage              - Run tests with coverage"
	@echo "  test-coverage-nocache      - Run tests with coverage without cache"
	@echo "  test-coverage-profile      - Run tests with coverage profile"
	@echo "  test-coverage-profile-nocache - Run tests with coverage profile without cache"
	@echo "  test-coverage-html         - Run tests with coverage and generate HTML report"
	@echo "  test-coverage-html-nocache - Run tests with coverage and generate HTML report without cache"
	@echo "  test-race                  - Run tests with race detection"
	@echo "  test-race-nocache          - Run tests with race detection without cache"
	@echo "  benchmark                  - Run benchmarks"
	@echo "  benchmark-nocache          - Run benchmarks without cache"
	@echo "  benchmark-race             - Run benchmarks with race detection"
	@echo "  benchmark-race-nocache     - Run benchmarks with race detection without cache"
	@echo "  tidy                       - Tidy go.mod"
	@echo "  install                    - Install dependencies"
	@echo "  clean                      - Clean build artifacts"
	@echo "  fmt                        - Format code"
	@echo "  vet                        - Vet code"
	@echo "  lint                       - Check for style issues using golangci-lint"
	@echo "  generate                   - Generate mocks"
	@echo "  check                      - Run all checks (vet, fmt, test)"
	@echo "  check-nocache              - Run all checks without cache"
	@echo "  check-race                 - Run all checks with race detection"
	@echo "  check-race-nocache         - Run all checks with race detection without cache"
	@echo "  install-tools              - Install development tools"
	@echo "  update                     - Update dependencies"
	@echo "  coverage                   - Show test coverage"
	@echo "  coverage-nocache           - Show test coverage without cache"
	@echo "  godoc                      - Generate documentation"
	@echo "  test-integration           - Run integration tests"
	@echo "  test-integration-nocache   - Run integration tests without cache"
	@echo "  test-all                   - Run all tests including integration tests"
	@echo "  test-all-nocache           - Run all tests including integration tests without cache"