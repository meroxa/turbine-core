module github.com/meroxa/turbine-core/v2

go 1.21.4

// Remove once conduit commons has been versioned
replace github.com/conduitio/conduit-commons v0.0.0-20231215130533-b393cff920e2 => github.com/conduitio/conduit-commons v0.0.0-20231222110339-d2dade3dc74a

require (
	github.com/conduitio/conduit-commons v0.0.0-20231215130533-b393cff920e2
	github.com/envoyproxy/protoc-gen-validate v1.0.2
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.5.0
	github.com/heimdalr/dag v1.4.0
	github.com/santhosh-tekuri/jsonschema/v5 v5.3.1
	github.com/stretchr/testify v1.8.4
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231211222908-989df2bf70f3 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
