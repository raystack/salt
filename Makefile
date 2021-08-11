.PHONY: check fmt lint test vet help
.DEFAULT_GOAL := help

check: test lint ## Run tests and linters

test: ## Run tests
	go test ./... -race

lint: ## Run linter
	golangci-lint run

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'