export CGO_ENABLED=0

# Target "all" should stay at the top of the file, we want it to be the default target.
.PHONY: all
all: vet test

### Help (lists all documented targets) ###

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+%?:.*?## .*$$' $(MAKEFILE_LIST) | sed 's/^*://g' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

### General targets ###

.PHONY: build
build:
	@echo "Nothing to build in turbine-core"

.PHONY: clean
clean:
	@echo "Nothing to clean in  turbine-core"

.PHONY: test
test: ## Run unit tests.
	CGO_ENABLED=1 go test $(GOTEST_FLAGS) -short -race -cover -covermode=atomic ./...

### Custom targets ###

.PHONY: fmt
fmt: ## Format Go files using gofumpt and gci.
	gofumpt -l -w .
	gci write --skip-generated  .

.PHONY: generate
generate: ## Run go generate.
	go generate ./...

.PHONY: vet
vet: ## Run go vet.
	go vet ./...

.PHONY: lint
lint: ## Lint Go files using golangci-lint.
	golangci-lint run -v

.PHONY: test-integration
test-integration: ## Run integration tests.
	go test $(GOTEST_FLAGS) -cover -covermode=atomic -run Integration ./...

.PHONY: tools
tools: ## Run make in "tools", optionally add "tools-[target]" to run a specific target.
	make -C tools

.PHONY: proto
proto: ## Generate Turbine GoLang gRPC bindings
	docker run \
		--rm \
		-v $(CURDIR)/proto:/defs \
		-v $(CURDIR)/lib/go:/out \
		namely/protoc-all  \
			-f ./turbine_v1.proto \
			-l go --with-validator -o /out
ruby-sdk-%:
	make -C $(CURDIR)/lib/ruby $*

