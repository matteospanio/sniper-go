GO = go
YARN = yarn
BINARY_NAME = sniper-go
PROD_PATH = /usr/local/share/$(BINARY_NAME)

.PHONY: help build setup dependencies install service clean
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
	python3 -c "$$PRINT_HELP_PYSCRIPT" < $(MAKEFILE_LIST)

clean: ## Remove build artifacts
	rm -rf dist
	rm -rf bin

build: ## Build the binary file
	@echo "Building the backend..."
	cd ./server && $(GO) build -o ../bin/$(BINARY_NAME) .

setup: ## Setup the project
	@echo "Setting up the project..."
	cd ./server && $(GO) get .
	@echo "Building the frontend..."
	$(YARN) cd ./client && install
	$(YARN) cd ./client && build:dev

dependencies: ## Install the dependencies
	@echo "Installing the dependencies..."
	bash ./scripts/install_dependencies.sh

install: dependencies setup build ## Install the project
	@echo "Installing the project..."
	mkdir -p $(PROD_PATH)
	cp -r README.md $(PROD_PATH)
	cp -r ./dist $(PROD_PATH)
	cp -r ./scripts $(PROD_PATH)
	cp -r ./templates $(PROD_PATH)
	cp ./bin/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "Done!"

service: install ## Create the systemd service
	@echo "Creating systemd service..."
	bash ./scripts/install_service.sh