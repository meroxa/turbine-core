package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	pb "github.com/meroxa/turbine-core/v2/lib/go/github.com/meroxa/turbine/core"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fixtureRecord struct {
	Key       interface{}
	Value     map[string]interface{}
	Timestamp string
}

func ReadFixture(ctx context.Context, file string) ([]*pb.Record, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var (
		rr             []*pb.Record
		fixtureRecords []fixtureRecord
	)

	if err := json.Unmarshal(b, &fixtureRecords); err != nil {
		return nil, err
	}

	for _, r := range fixtureRecords {
		rr = append(rr, wrapRecord(r))
	}

	return rr, nil
}

func wrapRecord(m fixtureRecord) *pb.Record {
	b, _ := json.Marshal(m.Value)

	ts := timestamppb.New(time.Now())
	if m.Timestamp != "" {
		t, _ := time.Parse(time.RFC3339, m.Timestamp)
		ts = timestamppb.New(t)
	}

	return &pb.Record{
		Key:       fmt.Sprintf("%v", m.Key),
		Value:     b,
		Timestamp: ts,
	}
}

func PrintRecords(name string, sr *pb.StreamRecords) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintf(w, "Destination %s\n", name)
	fmt.Fprintf(w, "----------------------\n")
	fmt.Fprintln(w, "index\trecord")
	fmt.Fprintln(w, "----\t----")
	for i, r := range sr.Records {
		fmt.Fprintf(w, "%d\t%s\n", i, string(r.Value))
		fmt.Fprintln(w, "----\t----")
	}
	fmt.Fprintf(w, "records written\t%d\n", len(sr.Records))
	w.Flush()
}
