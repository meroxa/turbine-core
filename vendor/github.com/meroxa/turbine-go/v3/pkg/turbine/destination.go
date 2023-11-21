package turbine

import (
	"context"
)

type Destination interface {
	Write(Records, string) error
	WriteWithContext(context.Context, Records, string) error
	WriteWithConfig(Records, string, ConnectionOptions) error
	WriteWithConfigWithContext(context.Context, Records, string, ConnectionOptions) error
}
