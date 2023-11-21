package build

import (
	"context"

	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/client"
)

// TODO: refactor
type source struct {
	*pb.Source
	*pb.Destination
	client.Client
}

func (r *source) Records(collection string, cfg sdk.ConnectionOptions) (sdk.Records, error) {
	return r.RecordsWithContext(context.Background(), collection, cfg)
}

func (r *source) RecordsWithContext(ctx context.Context, collection string, cfg sdk.ConnectionOptions) (sdk.Records, error) {
	c, err := r.ReadCollection(ctx, &pb.ReadCollectionRequest{
		Source:     r.Source,
		Collection: collection,
		Configs:    connectionOptions(cfg),
	})
	if err != nil {
		return sdk.Records{}, err
	}

	return collectionToRecords(c), nil
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
	_, err := r.WriteCollectionToDestination(ctx, &pb.WriteCollectionRequest{
		Destination:           r.Destination,
		SourceCollection:      recordsToCollection(rr),
		DestinationCollection: collection,
		Configs:               connectionOptions(cfg),
	})

	return err
}
