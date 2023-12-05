package build

import (
	"context"

	sdk "github.com/meroxa/turbine-go/v3/pkg/turbine"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core/v2"
	client "github.com/meroxa/turbine-core/pkg/client/v2"
)

type destination struct {
	id string
	c  client.Client
}

func (d *destination) Write(rr sdk.Records) error {
	return d.WriteWithContext(context.Background(), rr)
}

func (d *destination) WriteWithContext(ctx context.Context, rr sdk.Records) error {
	if _, err := d.c.WriteRecords(ctx, &pb.WriteRecordsRequest{
		DestinationID: d.id,
		StreamRecords: fromRecords(rr),
	}); err != nil {
		return err
	}

	return nil
}
