# Modem SMS Register Makefile
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BINARY_NAME = modem-sms-register
BUILD_DIR = bin
PKG = github.com/mike-bionic/modem-sms-register
LDFLAGS = -ldflags "-X main.version=${VERSION}"

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod

# Default target
.PHONY: all
all: clean build

# Build for current platform
.PHONY: build
build:
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd

# Build for all platforms
.PHONY: build-all
build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd
	
	# Linux 386
	GOOS=linux GOARCH=386 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-386 ./cmd
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd
	
	# Linux ARM
	GOOS=linux GOARCH=arm $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm ./cmd
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd
	
	# Windows 386
	GOOS=windows GOARCH=386 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-386.exe ./cmd
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd
	
	# macOS ARM64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd

# Build Windows GUI versions (no console window)
.PHONY: build-windows-gui
build-windows-gui:
	@echo "Building Windows GUI versions..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -ldflags "-H=windowsgui" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64-gui.exe ./cmd
	GOOS=windows GOARCH=386 $(GOBUILD) $(LDFLAGS) -ldflags "-H=windowsgui" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-386-gui.exe ./cmd

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Tidy dependencies
.PHONY: tidy
tidy:
	$(GOMOD) tidy

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download

# Format code
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Run the application with default config
.PHONY: run
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run with verbose logging
.PHONY: run-verbose
run-verbose: build
	./$(BUILD_DIR)/$(BINARY_NAME) -v

# Create example config
.PHONY: config
config:
	@if [ ! -f config.json ]; then \
		echo "Creating example config.json..."; \
		cp config.example.json config.json; \
		echo "Please edit config.json with your settings"; \
	else \
		echo "config.json already exists"; \
	fi

# Install development tools
.PHONY: install-tools
install-tools:
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker build
.PHONY: docker-build
docker-build:
	docker build -t $(BINARY_NAME):$(VERSION) .
	docker build -t $(BINARY_NAME):latest .

# Docker compose up
.PHONY: docker-up
docker-up:
	docker-compose up -d

# Docker compose down
.PHONY: docker-down
docker-down:
	docker-compose down

# Docker logs
.PHONY: docker-logs
docker-logs:
	docker-compose logs -f

# Show version
.PHONY: version
version:
	@echo $(VERSION)

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build for current platform"
	@echo "  build-all      - Build for all platforms"
	@echo "  build-windows-gui - Build Windows GUI versions"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  tidy           - Tidy dependencies"
	@echo "  deps           - Download dependencies"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  run            - Run application"
	@echo "  run-verbose    - Run with verbose logging"
	@echo "  config         - Create example config"
	@echo "  install-tools  - Install development tools"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-up      - Start with docker-compose"
	@echo "  docker-down    - Stop docker-compose"
	@echo "  docker-logs    - Show docker logs"
	@echo "  version        - Show version"
	@echo "  help           - Show this help"