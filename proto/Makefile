# Target "all" should stay at the top of the file, we want it to be the default target.
.PHONY: all
all: lint generate

### Help (lists all documented targets) ###

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+%?:.*?## .*$$' $(MAKEFILE_LIST) | sed 's/^*://g' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: generate
generate: ## Generate code for Turbine and Process with protoc plugins
	buf generate

.PHONY: update
update: ## Update a module's dependencies by updating the buf.lock file
	buf mod update

.PHONY: lint
lint: ## Run linting on Protobuf files
	buf lint

.PHONY: fmt
fmt: ## Format Protobuf files
	buf format -w