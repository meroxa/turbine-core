package build

import (
	"context"

	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/client"
)

// TODO: refactor
type destination struct {
	*pb.Source
	*pb.Destination
	client.Client
}

func (d *destination) Records(collection string, cfg sdk.ConnectionOptions) (sdk.Records, error) {
	return d.RecordsWithContext(context.Background(), collection, cfg)
}

func (d *destination) RecordsWithContext(ctx context.Context, collection string, cfg sdk.ConnectionOptions) (sdk.Records, error) {
	c, err := d.ReadCollection(ctx, &pb.ReadCollectionRequest{
		Source:     d.Source,
		Collection: collection,
		Configs:    connectionOptions(cfg),
	})
	if err != nil {
		return sdk.Records{}, err
	}

	return collectionToRecords(c), nil
}

func (d *destination) Write(rr sdk.Records, collection string) error {
	return d.WriteWithConfigWithContext(context.Background(), rr, collection, sdk.ConnectionOptions{})
}

func (d *destination) WriteWithContext(ctx context.Context, rr sdk.Records, collection string) error {
	return d.WriteWithConfigWithContext(ctx, rr, collection, sdk.ConnectionOptions{})
}

func (d *destination) WriteWithConfig(rr sdk.Records, collection string, cfg sdk.ConnectionOptions) error {
	return d.WriteWithConfigWithContext(context.Background(), rr, collection, cfg)
}

func (d *destination) WriteWithConfigWithContext(ctx context.Context, rr sdk.Records, collection string, cfg sdk.ConnectionOptions) error {
	_, err := d.WriteCollectionToDestination(ctx, &pb.WriteCollectionRequest{
		Destination:           d.Destination,
		SourceCollection:      recordsToCollection(rr),
		DestinationCollection: collection,
		Configs:               connectionOptions(cfg),
	})

	return err
}
