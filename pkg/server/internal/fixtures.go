package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fixtureRecord struct {
	Key       string
	Value     map[string]interface{}
	Timestamp string
}

func ReadFixtures(path, collection string) (*pb.Collection, error) {
	log.Printf("fixtures path: %s", path)
	b, err := os.ReadFile(path)
	if err != nil {
		return &pb.Collection{}, err
	}

	var records map[string][]fixtureRecord
	err = json.Unmarshal(b, &records)
	if err != nil {
		return &pb.Collection{}, err
	}

	var rr []*pb.Record
	for _, r := range records[collection] {
		rr = append(rr, wrapRecord(r))
	}

	col := &pb.Collection{
		Name:    collection,
		Records: rr,
	}

	return col, nil
}

func wrapRecord(m fixtureRecord) *pb.Record {
	b, _ := json.Marshal(m.Value)

	var ts *timestamppb.Timestamp
	if m.Timestamp == "" {
		ts = timestamppb.New(time.Now())
	} else {
		t, _ := time.Parse(time.RFC3339, m.Timestamp)
		ts = timestamppb.New(t)
	}

	return &pb.Record{
		Key:       m.Key,
		Value:     b,
		Timestamp: ts,
	}
}

func PrettyPrintRecords(name string, collection string, rr []*pb.Record) {
	fmt.Printf("=====================to %s (%s) resource=====================\n", name, collection)
	for _, r := range rr {
		payloadVal := string(r.Value)
		fmt.Println(payloadVal)
	}
	fmt.Printf("%d record(s) written\n", len(rr))
}
