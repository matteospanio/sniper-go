GO = go
YARN = yarn
BINARY_NAME = sniper-go

.PHONY: help build setup

help: ## Show this help
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  build: Build the binary file"
	@echo "  setup: Setup the project"

build: ## Build the binary file
	$(YARN) build:dev
	cd ./src && $(GO) build -o ../bin/$(BINARY_NAME) .

setup: ## Setup the project
	cd ./src && $(GO) get .
	$(YARN) install
