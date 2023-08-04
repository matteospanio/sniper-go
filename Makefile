GO = go
BINARY_NAME = sniper-go

.PHONY: help build setup

help: ## Show this help
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  build: Build the binary file"
	@echo "  setup: Setup the project"

build: src/main.go ## Build the binary file
	$(GO) build -o bin/$(BINARY_NAME) server/main.go

setup: ## Setup the project
	$(GO) get .
