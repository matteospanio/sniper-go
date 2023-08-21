GO := go
YARN := yarn
BINARY_NAME := sniper-go
SRC_DIR := server
BUILD_DIR := bin
PROD_PATH := /usr/local/share/$(BINARY_NAME)

.PHONY: help build setup dependencies install service clean docs
.DEFAULT_GOAL := help

define PRINT_HELP_PYSCRIPT
import re, sys

BOLD = '\033[1m'
BLUE = '\033[94m'
END = '\033[0m'

print("Usage: make <target>\n")
print(BOLD + "%-20s%s" % ("target", "description") + END)
for line in sys.stdin:
	match = re.match(r'^([a-zA-Z_-]+):.*?## (.*)$$', line)
	if match:
		target, help = match.groups()
		print( BLUE + "%-20s" % (target) + END + "%s" % (help))
endef
export PRINT_HELP_PYSCRIPT


help: ## Show this help
	@python3 -c "$$PRINT_HELP_PYSCRIPT" < $(MAKEFILE_LIST)

clean: ## Remove build artifacts
	rm -rf dist
	rm -rf $(BUILD_DIR)
	rm -rf site

build: ## Build the binary file
	@echo "Building the backend..."
	cd $(SRC_DIR) && $(GO) build -o ../$(BUILD_DIR)/$(BINARY_NAME) .
	$(YARN) --cwd ./client build:dev

setup: ## Setup the project
	@echo "Setting up the project..."
	cd $(SRC_DIR) && $(GO) get .
	@echo "Building the frontend..."
	$(YARN) --cwd ./client install

dependencies: ## Install the dependencies
	@echo "Installing the dependencies..."
	bash ./scripts/install_dependencies.sh

install: dependencies setup build ## Install the project
	@echo "Installing the project..."
	mkdir -p $(PROD_PATH)
	cp README.md $(PROD_PATH)
	cp -r ./dist $(PROD_PATH)
	cp -r ./scripts $(PROD_PATH)
	cp -r ./templates $(PROD_PATH)
	cp ./$(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "Done!"

service: install ## Create the systemd service
	@echo "Creating systemd service..."
	bash ./scripts/install_service.sh

docs: ## Generate the documentation
	@echo "Generating the documentation..."
	python3 -m mkdocs build --clean