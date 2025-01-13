# Colors for pretty output
.PHONY: build test test-v clean fmt lint help build-all install

GREEN  := $(shell tput -Txterm setaf 2)
RESET  := $(shell tput -Txterm sgr0)
BOLD   := $(shell tput -Txterm bold)

define success
	echo "${GREEN}${BOLD}✓ $1${RESET}"
endef

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "development")

# Binary Info
BINARY_NAME = seed
BINARY_DIR = bin
ifeq ($(OS),Windows_NT)
    BINARY_EXTENSION = .exe
else
    BINARY_EXTENSION =
endif

# Main package path reflecting the new structure
MAIN_PATH = cmd/seed

# Build flags
LDFLAGS = -ldflags "-X main.version=$(VERSION)"

# Default target when just running 'make'
all: install-deps

# Success message function

# Build the binary
build:
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	@mkdir -p $(BINARY_DIR)
	go build $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)$(BINARY_EXTENSION) ./$(MAIN_PATH)
	@$(call success,"Built")


build-all:
	@echo "Building $(BINARY_NAME) for all platforms..."
	@mkdir -p $(BINARY_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(MAIN_PATH)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-linux-amd64 ./$(MAIN_PATH)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(MAIN_PATH)
	@$(call success,"Built all")

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
	@echo "Cleaning build artifacts..."
	go clean
	rm -rf $(BINARY_DIR)
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

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@install -m 755 $(BINARY_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@$(call success,"Installed. Plant some seeds!")

# Install development dependencies
install-deps:
	@echo "Installing dependencies..."
	@go mod download
	@$(call success,"All dependencies installed")
	# Then install development tools
	@echo "Installing development tools..."
	@go install gotest.tools/gotestsum@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	
	@$(call success,"Installing development tools. Ready to develop!")

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
