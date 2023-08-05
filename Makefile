GO = go
BINARY_NAME = sniper-go

.PHONY: help build setup

help: ## Show this help
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  build: Build the binary file"
	@echo "  setup: Setup the project"

build: ## Build the binary file
	$(GO) build -o bin/$(BINARY_NAME) src/main.go

setup: ## Setup the project
	$(GO) get .
