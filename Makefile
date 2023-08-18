GO = go
YARN = yarn
BINARY_NAME = sniper-go

.PHONY: help build setup dependencies install service clean
.DEFAULT_GOAL := help

clean: ## Remove build files
	rm -rf dist
    rm -rf bin

help: ## Show this help
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  install: Install the project (requires sudo)"
	@echo "  dependencies: Install the dependencies"
	@echo "  setup: Setup the project"
	@echo "  build: Build the binary file"
	@echo "  service: Create the systemd service"

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
	@echo "Installing the project..."
	sudo mkdir -p /usr/local/share/$(BINARY_NAME)
	sudo cp -r README.md /usr/local/share/$(BINARY_NAME)
	sudo cp -r ./dist /usr/local/share/$(BINARY_NAME)
	sudo cp -r ./scripts /usr/local/share/$(BINARY_NAME)
	sudo cp -r ./templates /usr/local/share/$(BINARY_NAME)
	sudo cp ./bin/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "Done!"

service: install ## Create the systemd service
	@echo "Creating systemd service..."
	bash ./scripts/install_service.sh