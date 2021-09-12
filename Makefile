# This is main makefile
.DEFAULT_GOAL := help
.PHONY:

WORKSPACE := $(shell pwd)

#Print help message
help: ## Get help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

codegen: ## Generate code from swagger
	@echo "Generate code from swagger..."
	@swagger generate server -A user-store -t internal --exclude-main -f ./swagger/user-store.yaml
	@echo "Add dependency modules..."
	@go mod tidy -v

units: codegen ## Run Unit Tests
	@echo "Run Unit Tests..."
	@go test ./... -cover

build: units ## Build container image
	@echo "Build container image..."
	@docker build -t ilyatg/userstore:1.0.0 .

run: build ## Run container image
	@echo "Run container image..."
	@docker run --rm -p 9090:9090 -e PORT=9090 ilyatg/userstore:1.0.0