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
	"google.golang.org/protobuf/types/known/structpb"
)

func ToProto(records []opencdc.Record) ([]*procproto.Record, error) {
	out := make([]*procproto.Record, len(records))
	for i, record := range records {
		outRec, err := protoRecord(record)
		if err != nil {
			return nil, fmt.Errorf("failed converting record %v to proto: %w", i, err)
		}
		out[i] = outRec
	}

	return out, nil
}

func protoRecord(record opencdc.Record) (*procproto.Record, error) {
	key, err := protoData(record.Key)
	if err != nil {
		return nil, fmt.Errorf("error converting key: %w", err)
	}

	payload, err := protoChange(record.Payload)
	if err != nil {
		return nil, fmt.Errorf("error converting payload: %w", err)
	}

	out := procproto.Record{
		Position:  record.Position,
		Operation: procproto.Operation(record.Operation),
		Metadata:  record.Metadata,
		Key:       key,
		Payload:   payload,
	}
	return &out, nil
}

func protoChange(in opencdc.Change) (*procproto.Change, error) {
	before, err := protoData(in.Before)
	if err != nil {
		return nil, fmt.Errorf("error converting before: %w", err)
	}

	after, err := protoData(in.After)
	if err != nil {
		return nil, fmt.Errorf("error converting after: %w", err)
	}

	out := procproto.Change{
		Before: before,
		After:  after,
	}
	return &out, nil
}

func protoData(in opencdc.Data) (*procproto.Data, error) {
	if in == nil {
		return nil, nil
	}

	var out procproto.Data

	switch v := in.(type) {
	case opencdc.RawData:
		out.Data = &procproto.Data_RawData{
			RawData: v,
		}
	case opencdc.StructuredData:
		content, err := structpb.NewStruct(v)
		if err != nil {
			return nil, err
		}
		out.Data = &procproto.Data_StructuredData{
			StructuredData: content,
		}
	default:
		return nil, errors.New("invalid Data type")
	}

	return &out, nil
}
