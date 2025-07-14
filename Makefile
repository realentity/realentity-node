# Makefile for RealEntity Node

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# Docker variables
DOCKER_IMAGE = realentity/node
DOCKER_TAG ?= latest

# Go variables
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= 0

# Directories
BUILD_DIR = build
DEPLOY_DIR = deploy
DOCS_DIR = docs

.PHONY: help build test clean deploy docker lint fmt vet mod-tidy

help: ## Show this help message
	@echo "RealEntity Node Build System"
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building RealEntity Node..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/realentity-node cmd/main.go
	@echo "Build complete: $(BUILD_DIR)/realentity-node"

build-all: ## Build for all platforms
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			if [ "$$os" = "windows" ]; then \
				ext=".exe"; \
			else \
				ext=""; \
			fi; \
			echo "Building $$os/$$arch..."; \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build $(LDFLAGS) -o $(BUILD_DIR)/realentity-node-$$os-$$arch$$ext cmd/main.go; \
		done \
	done
	@echo "Multi-platform build complete"

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	go test -v -tags=integration ./test/integration/...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linters
	@echo "Running linters..."
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

mod-tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	go mod tidy

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -f realentity-node realentity-node.exe

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_IMAGE):$(VERSION)

docker-push: ## Push Docker image
	@echo "Pushing Docker image..."
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_IMAGE):$(VERSION)

docker-test: ## Run Docker test environment
	@echo "Starting Docker test environment..."
	./$(DEPLOY_DIR)/universal.sh docker

deploy-local: ## Deploy for local development
	@echo "Deploying for local development..."
	./$(DEPLOY_DIR)/universal.sh local

deploy-bootstrap: ## Deploy VPS bootstrap node
	@echo "Deploying VPS bootstrap node..."
	@read -p "Enter public IP: " public_ip; \
	./$(DEPLOY_DIR)/universal.sh vps-bootstrap --public-ip $$public_ip

deploy-peer: ## Deploy VPS peer node
	@echo "Deploying VPS peer node..."
	@read -p "Enter bootstrap peer address: " bootstrap_peer; \
	./$(DEPLOY_DIR)/universal.sh vps-peer --bootstrap-peer "$$bootstrap_peer"

install-deps: ## Install development dependencies
	@echo "Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/godoc@latest

release: clean test build-all docker-build ## Create a release
	@echo "Creating release $(VERSION)..."
	@mkdir -p $(BUILD_DIR)/release
	@cp $(BUILD_DIR)/realentity-node-* $(BUILD_DIR)/release/
	@cd $(BUILD_DIR)/release && tar -czf realentity-node-$(VERSION).tar.gz realentity-node-*
	@echo "Release created: $(BUILD_DIR)/release/realentity-node-$(VERSION).tar.gz"

validate-config: ## Validate configuration files
	@echo "Validating configuration files..."
	@for config in config*.json configs/*.json $(DEPLOY_DIR)/configs/*.json; do \
		if [ -f "$$config" ]; then \
			echo "Validating $$config..."; \
			python -m json.tool "$$config" > /dev/null && echo " $$config is valid" || echo " $$config is invalid"; \
		fi \
	done

migrate: ## Migrate from legacy deployment scripts
	@echo "Migrating deployment scripts..."
	chmod +x migrate-deployment.sh
	./migrate-deployment.sh interactive

clean-legacy: ## Remove legacy deployment scripts (after migration)
	@echo "Cleaning up legacy files..."
	@if [ -f "migrate-deployment.sh" ]; then \
		./migrate-deployment.sh migrate; \
	else \
		echo "Run 'make migrate' first"; \
	fi

validate-deployment: ## Validate deployment configuration
	@echo "Validating deployment setup..."
	@if [ ! -f "deploy/universal.sh" ]; then \
		echo " deploy/universal.sh not found"; \
		exit 1; \
	fi
	@if [ ! -x "deploy/universal.sh" ]; then \
		echo " deploy/universal.sh not executable"; \
		exit 1; \
	fi
	@echo " Deployment system validated"

setup-dev: install-deps validate-deployment ## Setup development environment
	@echo "Setting up development environment..."
	@if [ ! -f config.json ]; then \
		echo "Creating local development config..."; \
		./$(DEPLOY_DIR)/universal.sh local --dry-run; \
	fi
	@echo "Development environment ready!"

# CI/CD targets
ci-test: mod-tidy fmt vet lint test ## Run all CI tests
	@echo "All CI tests passed!"

ci-build: ci-test build docker-build ## Run CI build pipeline
	@echo "CI build completed!"

# Help is default target
.DEFAULT_GOAL := help
