package tools

import (
	//nolint:typecheck // this is a list of tools
	_ "github.com/daixiang0/gci"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/vektra/mockery/v2"
	_ "mvdan.cc/gofumpt"
)
