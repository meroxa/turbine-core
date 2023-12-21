// Copyright © 2022 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package record

import (
	"errors"
	"fmt"

	"github.com/conduitio/conduit-commons/opencdc"
	procproto "github.com/conduitio/conduit-connector-protocol/proto/opencdc/v1"
)

func FromProto(in []*procproto.Record) ([]opencdc.Record, error) {
	if in == nil {
		return nil, nil
	}

	outRecs := make([]opencdc.Record, len(in))
	for i, protoRec := range in {
		rec, err := opencdcRecord(protoRec)
		if err != nil {
			return nil, fmt.Errorf("failed converting protobuf record at index %v to OpenCDC record: %w", i, err)
		}

		outRecs[i] = rec
	}

	return outRecs, nil
}

func opencdcRecord(record *procproto.Record) (opencdc.Record, error) {
	key, err := opencdcData(record.Key)
	if err != nil {
		return opencdc.Record{}, fmt.Errorf("error converting key: %w", err)
	}

	payload, err := opencdcChange(record.Payload)
	if err != nil {
		return opencdc.Record{}, fmt.Errorf("error converting payload: %w", err)
	}

	out := opencdc.Record{
		Position:  record.Position,
		Operation: opencdc.Operation(record.Operation),
		Metadata:  record.Metadata,
		Key:       key,
		Payload:   payload,
	}
	return out, nil
}

func opencdcChange(in *procproto.Change) (opencdc.Change, error) {
	before, err := opencdcData(in.Before)
	if err != nil {
		return opencdc.Change{}, fmt.Errorf("error converting before: %w", err)
	}

	after, err := opencdcData(in.After)
	if err != nil {
		return opencdc.Change{}, fmt.Errorf("error converting after: %w", err)
	}

	out := opencdc.Change{
		Before: before,
		After:  after,
	}
	return out, nil
}

func opencdcData(in *procproto.Data) (opencdc.Data, error) {
	d := in.GetData()
	if d == nil {
		return nil, nil
	}

	switch v := d.(type) {
	case *procproto.Data_RawData:
		return opencdc.RawData(v.RawData), nil
	case *procproto.Data_StructuredData:
		return opencdc.StructuredData(v.StructuredData.AsMap()), nil
	default:
		return nil, errors.New("invalid Data type")
	}
}
