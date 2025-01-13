# Colors for pretty output
GREEN  := $(shell tput -Txterm setaf 2)
RESET  := $(shell tput -Txterm sgr0)
BOLD   := $(shell tput -Txterm bold)

.PHONY: build test test-v clean fmt lint help

# Default target when just running 'make'
all: build

# Success message function
define success
	echo "${GREEN}${BOLD}✓ $1${RESET}"
endef

# Build the binary
build:
	@echo "Building..."
	@go build -o bin/seed
	@$(call success,"Built")

# Run tests
test:
	@gotestsum --format pkgname-and-test-fails -- ./... && $(call success,"All tests passed")

# Run tests with verbose output
test-v:
	@gotestsum --format standard-verbose -- -v ./... && $(call success,"All tests passed")

# Run tests and watch them
test-w:
	@gotestsum --watch --format standard-verbose -- -v ./... 

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f bin/seed
	@go clean
	@$(call success,"Sparkling fresh and new")

# Format code
fmt:
	@echo "Formatting..."
	@go fmt ./...
	@$(call success,"Formatting completed")

# Run linter
lint:
	@echo "Linting..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run && $(call success,"All looks good!"); \
	else \
		echo "golangci-lint is not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

# Install development dependencies
install-deps:
	@echo "Installing development dependencies..."
	@go install gotest.tools/gotestsum@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(call success,"Deps installed. Lets rock and roll.")

# Show help
help:
	@echo "Available targets:"
	@echo "  make          - Build the binary"
	@echo "  make build    - Build the binary"
	@echo "  make test     - Run tests"
	@echo "  make test-v   - Run tests with verbose output"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make fmt      - Format code"
	@echo "  make lint     - Run linter"
	@echo "  make install-deps - Install development dependencies"
	@echo "  make help     - Show this help message"
