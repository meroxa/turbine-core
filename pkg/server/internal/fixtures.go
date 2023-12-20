package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/conduitio/conduit-commons/opencdc"
	"github.com/conduitio/conduit-connector-protocol/proto/opencdc/v1"
	"github.com/meroxa/turbine-core/v2/pkg/record"
	pb "github.com/meroxa/turbine-core/v2/proto/turbine/v2"
)

func ReadFixture(ctx context.Context, file string) ([]*opencdcv1.Record, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var fixtureRecords []opencdc.Record

	if err := json.Unmarshal(b, &fixtureRecords); err != nil {
		return nil, err
	}

	rr, err := record.ToProto(fixtureRecords)
	if err != nil {
		return nil, err
	}

	return rr, nil
}

func PrintRecords(name string, sr *pb.StreamRecords) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintf(w, "Destination %s\n", name)
	fmt.Fprintf(w, "----------------------\n")
	fmt.Fprintln(w, "index\trecord")
	fmt.Fprintln(w, "----\t----")

	rr, err := record.FromProto(sr.Records)
	if err != nil {
		panic(err)
	}

	for i, r := range rr {
		fmt.Fprintf(w, "%d\t%s\n", i, string(r.Bytes()))
		fmt.Fprintln(w, "----\t----")
	}

	fmt.Fprintf(w, "records written\t%d\n", len(sr.Records))
	w.Flush()
}
