MOCKGEN_VERSION   ?= v1.6.0

.PHONY: gomod
gomod:
	go mod tidy && go mod vendor

.PHONY: vet
vet:
	go vet ./...

.PHONY: gofumpt
gofumpt:
	go install mvdan.cc/gofumpt@latest

.PHONY: fmt
fmt: gofumpt
	gofumpt -l -w .

.PHONY: generate
generate: mockgen-install
	go generate ./...

.PHONY: test
test: vet
	go test `go list ./... | grep -v 'turbine-core\/pkg\/app\/templates\/go'` \
	-coverprofile=c.out \
	-covermode=atomic \
	-v

.PHONY: test_turbine_rb
test_turbine_rb:
	cd lib/ruby/turbine_rb && \
		bundler install && \
		bundler exec rake

.PHONY: proto
proto: turbine_proto process_ruby_proto turbine_ruby_proto

.PHONY: turbine_proto
turbine_proto:
	docker run \
		--rm \
		-v $(CURDIR)/proto:/defs \
		-v $(CURDIR)/lib/go:/out \
		namely/protoc-all  \
			-f ./turbine/v1/turbine.proto \
			-l go --with-validator -o /out

.PHONY: process_ruby_proto
process_ruby_proto:
	docker run \
		--rm \
		-v $(CURDIR)/proto/process/v1:/defs \
		-v $(CURDIR)/lib/ruby/turbine_rb/lib:/out \
		namely/protoc-all  \
			-f ./service.proto \
			-l ruby -o /out

.PHONY: turbine_ruby_proto
turbine_ruby_proto:
	docker run \
		--rm \
		-v $(CURDIR)/proto/turbine/v1:/defs \
		-v $(CURDIR)/lib/ruby/turbine_rb/lib:/out \
		namely/protoc-all  \
			-f ./turbine.proto \
			-l ruby -o /out

mockgen-install:
	go install github.com/golang/mock/mockgen@$(MOCKGEN_VERSION)