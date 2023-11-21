package server

import (
	"context"

	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"
)

var _ sdk.Source = (*source)(nil)

type source struct{}

func (r *source) Records(collection string, cfg sdk.ConnectionOptions) (sdk.Records, error) {
	return r.RecordsWithContext(context.Background(), collection, cfg)
}

func (r *source) RecordsWithContext(ctx context.Context, collection string, cfg sdk.ConnectionOptions) (sdk.Records, error) {
	return sdk.Records{}, nil
}

func (r *source) Write(rr sdk.Records, collection string) error {
	return r.WriteWithConfigWithContext(context.Background(), rr, collection, sdk.ConnectionOptions{})
}

func (r *source) WriteWithContext(ctx context.Context, rr sdk.Records, collection string) error {
	return r.WriteWithConfigWithContext(ctx, rr, collection, sdk.ConnectionOptions{})
}

func (r *source) WriteWithConfig(rr sdk.Records, collection string, cfg sdk.ConnectionOptions) error {
	return r.WriteWithConfigWithContext(context.Background(), rr, collection, cfg)
}

func (r *source) WriteWithConfigWithContext(ctx context.Context, rr sdk.Records, collection string, cfg sdk.ConnectionOptions) error {
	return nil
}
