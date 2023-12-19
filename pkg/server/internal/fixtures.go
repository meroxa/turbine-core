package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/conduitio/conduit-commons/opencdc"
	"github.com/conduitio/conduit-connector-protocol/proto/opencdc/v1"
	pb "github.com/meroxa/turbine-core/v2/proto/turbine/v2"
)

func ReadFixture(ctx context.Context, file string) ([]*opencdcv1.Record, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var (
		rr             []*opencdcv1.Record
		fixtureRecords []*opencdc.Record
	)

	if err := json.Unmarshal(b, &fixtureRecords); err != nil {
		return nil, err
	}

	for _, r := range fixtureRecords {
		rr = append(rr, wrapRecord(r))
	}

	return rr, nil
}

func wrapRecord(r *opencdc.Record) *opencdcv1.Record {
	return &opencdcv1.Record{
		Position:  []byte(r.Position),
		Operation: opencdcv1.Operation(r.Operation),
		Metadata:  r.Metadata,
		Key:       nil, /* &v1.Data{Data: raw/structured} */
		Payload:   nil, /* &v1.Change{Before: &v1.Data{Data: //}, After...} */
	}
}

func unwrapRecord(r *opencdcv1.Record) *opencdc.Record {
	return &opencdc.Record{}
}

func PrintRecords(name string, sr *pb.StreamRecords) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintf(w, "Destination %s\n", name)
	fmt.Fprintf(w, "----------------------\n")
	fmt.Fprintln(w, "index\trecord")
	fmt.Fprintln(w, "----\t----")
	for i, r := range sr.Records {
		unwrapped := unwrapRecord(r)
		fmt.Fprintf(w, "%d\t%s\n", i, string(unwrapped.Bytes()))
		fmt.Fprintln(w, "----\t----")
	}
	fmt.Fprintf(w, "records written\t%d\n", len(sr.Records))
	w.Flush()
}
