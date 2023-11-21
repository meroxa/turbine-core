package turbine

import (
	"context"
)

type Source interface {
	Records(string, ConnectionOptions) (Records, error)
	RecordsWithContext(context.Context, string, ConnectionOptions) (Records, error)
}
