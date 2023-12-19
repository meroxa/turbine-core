# Target "all" should stay at the top of the file, we want it to be the default target.
.PHONY: all
all: install

### Help (lists all documented targets) ###

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+%?:.*?## .*$$' $(MAKEFILE_LIST) | sed 's/^*://g' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

### General targets ###

.PHONY: build
build:
	@echo "Nothing to build in tools."

.PHONY: clean
clean:
	@echo "Nothing to clean in tools."

.PHONY: test
test:
	@echo "Nothing to test in tools."

### Custom targets ###

.PHONY: install
install: install-go ## Install all tools (Go).

.PHONY: install-go
install-go: ## Install all Go tools.
	go list -e -f '{{ join .Imports "\n" }}' tools.go | xargs -tI % go install %
