GO = go
YARN = yarn
BINARY_NAME = sniper-go

.PHONY: help build setup dependencies install

help: ## Show this help
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  install: Install the project (requires sudo)"
	@echo "  dependencies: Install the dependencies"
	@echo "  setup: Setup the project"
	@echo "  build: Build the binary file"

build: ## Build the binary file
	@echo "Building the backend..."
	cd ./src && $(GO) build -o ../bin/$(BINARY_NAME) .

setup: ## Setup the project
	@echo "Setting up the project..."
	cd ./src && $(GO) get .
	$(YARN) install
	@echo "Building the frontend..."
	$(YARN) build:dev

dependencies: ## Install the dependencies
	@echo "Installing the dependencies..."
	bash ./scripts/install_dependencies.sh

install: dependencies setup build ## Install the project
	@echo "Done!"
