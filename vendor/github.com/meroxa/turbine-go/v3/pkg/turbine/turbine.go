package turbine

import (
	"context"
)

type Turbine interface {
	Source(string) (Source, error)
	SourceWithContext(context.Context, string) (Source, error)
	Destination(string) (Destination, error)
	DestinationWithContext(context.Context, string) (Destination, error)

	Process(Records, Function) (Records, error)
	ProcessWithContext(context.Context, Records, Function) (Records, error)

	RegisterSecret(name string) error
	RegisterSecretWithContext(ctx context.Context, name string) error
}
