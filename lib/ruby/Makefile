PROTO := $(abspath $(CURDIR)/../../proto)

all: test

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+%?:.*?## .*$$' $(MAKEFILE_LIST) | sed 's/^*://g' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Run ruby sdk tests
	cd ./turbine_rb && \
		bundler install && \
		bundler exec rake

.PHONY: proto
proto: process_proto turbine_proto # Generate turbine and function gRPC bindings

.PHONY: process_proto
process_proto: ## Generate function gRPC bindings
	docker run \
		--rm \
		-v $(PROTO):/defs \
		-v $(CURDIR)/turbine_rb/lib:/out \
		namely/protoc-all  \
			-f ./process_v1.proto \
			-l ruby -o /out

.PHONY: turbine_proto
turbine_proto: ## Generate turbine gRPC bindings
	docker run \
		--rm \
		-v $(PROTO):/defs \
		-v $(CURDIR)/turbine_rb/lib:/out \
		namely/protoc-all  \
			-f ./turbine_v2.proto \
			-l ruby -o /out

