package server

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/conduitio/conduit-connector-protocol/proto/opencdc/v1"
	"github.com/meroxa/turbine-core/v2/pkg/app"
	"github.com/meroxa/turbine-core/v2/pkg/ir"
	"github.com/meroxa/turbine-core/v2/proto/turbine/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

//go:embed testdata/opencdc_record.json
var testOpenCDCRecord []byte

func Test_Init(t *testing.T) {
	ctx := context.Background()
	tempdir := t.TempDir()
	tests := []struct {
		desc    string
		setup   func() *turbinev2.InitRequest
		wantErr error
	}{
		{
			desc:    "fails with invalid app name",
			wantErr: errors.New("invalid InitRequest.AppName: value length must be at least 1 runes"),
			setup: func() *turbinev2.InitRequest {
				return &turbinev2.InitRequest{
					ConfigFilePath: "/foo/bar",
					Language:       turbinev2.Language_GOLANG,
				}
			},
		},
		{
			desc:    "fails with invalid config file path",
			wantErr: errors.New("invalid InitRequest.ConfigFilePath: value length must be at least 1 runes"),
			setup: func() *turbinev2.InitRequest {
				return &turbinev2.InitRequest{
					AppName:  "turbine-app",
					Language: turbinev2.Language_GOLANG,
				}
			},
		},
		{
			desc:    "fails with invalid lang",
			wantErr: errors.New("invalid InitRequest.Language: value must be one of the defined enum values"),
			setup: func() *turbinev2.InitRequest {
				return &turbinev2.InitRequest{
					AppName:        "turbine-app",
					ConfigFilePath: "/foo/bar",
					Language:       101221,
				}
			},
		},
		{
			desc:    "fails to load app config",
			wantErr: errors.New("no such file or directory"),
			setup: func() *turbinev2.InitRequest {
				return &turbinev2.InitRequest{
					AppName:        "test-app",
					ConfigFilePath: "/nonexistingapp",
				}
			},
		},
		{
			desc: "success",
			setup: func() *turbinev2.InitRequest {
				file := path.Join(tempdir, "app.json")
				require.NoError(
					t,
					os.WriteFile(
						file,
						[]byte(fmt.Sprintf(`{
							"name": "app",
							"language": "%s",
							"fixtures": {
								"demopg": "%s"
							}
						}`, ir.Ruby, filepath.Join("fixtures", "demo.json"))),
						0o644,
					),
				)

				return &turbinev2.InitRequest{
					AppName:        "app",
					ConfigFilePath: tempdir,
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			req := tc.setup()
			s := &RunService{}

			_, err := s.Init(ctx, req)
			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, s.appPath, req.ConfigFilePath)
				assert.Equal(t, s.config, app.Config{
					Name: "app",
					Fixtures: map[string]string{
						"demopg": filepath.Join("fixtures", "demo.json"),
					},
					Language: ir.Ruby,
				})
			}
		})
	}
}

func Test_AddSource(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		desc    string
		setup   func() *turbinev2.AddSourceRequest
		wantErr error
	}{
		{
			desc: "fails on invalid name",
			setup: func() *turbinev2.AddSourceRequest {
				return &turbinev2.AddSourceRequest{}
			},
			wantErr: errors.New("invalid AddSourceRequest.Name: value length must be at least 1 runes"),
		},
		{
			desc: "success",
			setup: func() *turbinev2.AddSourceRequest {
				return &turbinev2.AddSourceRequest{
					Name: "source-name",
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			s := &RunService{}
			req := tc.setup()

			r, err := s.AddSource(ctx, req)
			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else if assert.NoError(t, err) {
				assert.Equal(t, r.StreamName, req.Name)
			}
		})
	}
}

func Test_ReadRecords(t *testing.T) {
	ctx := context.Background()
	tempdir := t.TempDir()
	tests := []struct {
		desc        string
		srv         *RunService
		setup       func() *turbinev2.ReadRecordsRequest
		wantRecords *turbinev2.ReadRecordsResponse
		wantErr     error
	}{
		{
			desc:    "fails when source is missing",
			srv:     &RunService{},
			wantErr: errors.New("invalid ReadRecordsRequest.SourceStream: value length must be at least 1 runes"),
			setup: func() *turbinev2.ReadRecordsRequest {
				return &turbinev2.ReadRecordsRequest{}
			},
		},
		{
			desc: "fails on missing fixture file",
			srv: &RunService{
				appPath: tempdir,
				config: app.Config{
					Fixtures: map[string]string{
						"resource": "fixture.json",
					},
				},
			},
			wantErr: errors.New("no such file or directory"),
			setup: func() *turbinev2.ReadRecordsRequest {
				return &turbinev2.ReadRecordsRequest{
					SourceStream: "resource",
				}
			},
		},
		/* Skip until fixture serialization to Record works.
		{
			desc: "success",
			srv: &RunService{
				appPath: path.Join(tempdir),
				config: app.Config{
					Fixtures: map[string]string{
						"source": "fixture.json",
					},
				},
			},
			wantRecords: &turbinev2.ReadRecordsResponse{
				StreamRecords: &turbinev2.StreamRecords{
					StreamName: "source",
					Records:    []*opencdcv1.Record{testProtoRecord(t)},
				},
			},
			setup: func() *turbinev2.ReadRecordsRequest {
				file := path.Join(tempdir, "fixture.json")
				require.NoError(t, os.WriteFile(file, testOpenCDCRecord, 0o644))

				return &turbinev2.ReadRecordsRequest{SourceStream: "source"}
			},
		},
		*/
		{
			desc: "wrong fixture source name",
			srv: &RunService{
				appPath: path.Join(tempdir),
				config: app.Config{
					Fixtures: map[string]string{
						"source123": "fixture.json",
					},
				},
			},
			wantErr: errors.New("no fixture file found for source pg"),
			setup: func() *turbinev2.ReadRecordsRequest {
				file := path.Join(tempdir, "fixture.json")
				require.NoError(t, os.WriteFile(file, testOpenCDCRecord, 0o644))

				return &turbinev2.ReadRecordsRequest{SourceStream: "pg"}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			req := tc.setup()

			c, err := tc.srv.ReadRecords(ctx, req)
			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else if assert.NoError(t, err) {
				assert.Equal(t, c.StreamRecords.StreamName, tc.wantRecords.StreamRecords.StreamName)
				assert.Equal(t, len(c.StreamRecords.Records), len(tc.wantRecords.StreamRecords.Records))
				// assert.Equal(t, c.StreamRecords.Records[0], // opencdc record)
			}
		})
	}
}

func Test_AddDestination(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		desc    string
		setup   func() *turbinev2.AddDestinationRequest
		wantErr error
	}{
		{
			desc: "fails on invalid name",
			setup: func() *turbinev2.AddDestinationRequest {
				return &turbinev2.AddDestinationRequest{}
			},
			wantErr: errors.New("invalid AddDestinationRequest.Name: value length must be at least 1 runes"),
		},
		{
			desc: "success",
			setup: func() *turbinev2.AddDestinationRequest {
				return &turbinev2.AddDestinationRequest{
					Name: "destination-name",
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			s := &RunService{}
			req := tc.setup()

			r, err := s.AddDestination(ctx, req)
			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else if assert.NoError(t, err) {
				assert.Equal(t, r.StreamName, req.Name)
			}
		})
	}
}

func Test_WriteRecords(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		desc    string
		setup   func(*testing.T) *turbinev2.WriteRecordsRequest
		wantErr error
	}{
		{
			desc:    "fails when destinationID is missing",
			wantErr: errors.New("invalid WriteRecordsRequest.DestinationID: value length must be at least 1 runes"),
			setup: func(_ *testing.T) *turbinev2.WriteRecordsRequest {
				return &turbinev2.WriteRecordsRequest{}
			},
		},
		{
			desc:    "fails when streamRecords is missing",
			wantErr: errors.New("invalid WriteRecordsRequest.StreamRecords: value is required"),
			setup: func(_ *testing.T) *turbinev2.WriteRecordsRequest {
				return &turbinev2.WriteRecordsRequest{
					DestinationID: "stream-destination",
				}
			},
		},
		{
			desc: "success",
			setup: func(t *testing.T) *turbinev2.WriteRecordsRequest {
				t.Helper()

				return &turbinev2.WriteRecordsRequest{
					DestinationID: "destination-stream",
					StreamRecords: &turbinev2.StreamRecords{
						StreamName: "source",
						Records:    []*opencdcv1.Record{testProtoRecord(t)},
					},
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			s := &RunService{}
			req := tc.setup(t)

			// capture stdout and match if it contains what we need
			capture := func(fn func() error) (string, error) {
				return "", fn()
			}
			if tc.wantErr == nil {
				capture = func(fn func() error) (string, error) {
					stdout := os.Stdout
					r, w, err := os.Pipe()
					require.NoError(t, err)

					os.Stdout = w
					err = fn()
					w.Close()
					os.Stdout = stdout

					var buf bytes.Buffer
					io.Copy(&buf, r)
					return buf.String(), err
				}
			}

			output, err := capture(func() error {
				_, err := s.WriteRecords(ctx, req)
				return err
			})
			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				assert.Contains(t, output, string(testJSONRecord(t)))
				assert.Contains(t, output, "destination-stream")
				assert.NoError(t, err)
			}
		})
	}
}

func Test_ProcessRecords(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		desc    string
		setup   func() *turbinev2.ProcessRecordsRequest
		wantErr error
	}{
		{
			desc: "fails on missing process",
			setup: func() *turbinev2.ProcessRecordsRequest {
				return &turbinev2.ProcessRecordsRequest{}
			},
			wantErr: errors.New("invalid ProcessRecordsRequest.Process: value is required"),
		},
		{
			desc: "fails on missing streamRecords",
			setup: func() *turbinev2.ProcessRecordsRequest {
				return &turbinev2.ProcessRecordsRequest{
					Process: &turbinev2.ProcessRecordsRequest_Process{Name: "my-process"},
				}
			},
			wantErr: errors.New("invalid ProcessRecordsRequest.StreamRecords: value is required"),
		},
		{
			desc: "success",
			setup: func() *turbinev2.ProcessRecordsRequest {
				return &turbinev2.ProcessRecordsRequest{
					Process: &turbinev2.ProcessRecordsRequest_Process{Name: "my-process"},
					StreamRecords: &turbinev2.StreamRecords{
						StreamName: "my-stream",
					},
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			s := &RunService{}
			req := tc.setup()

			c, err := s.ProcessRecords(ctx, req)
			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else if assert.NoError(t, err) {
				assert.Equal(t, c.StreamRecords.StreamName, req.StreamRecords.StreamName)
			}
		})
	}
}

func testJSONRecord(t *testing.T) []byte {
	t.Helper()
	var out bytes.Buffer
	require.NoError(t, json.Compact(&out, testOpenCDCRecord))

	return out.Bytes()
}

func testProtoRecord(t *testing.T) *opencdcv1.Record {
	t.Helper()

	keydata, err := structpb.NewStruct(map[string]any{
		"id": 1,
	})
	require.NoError(t, err)

	afterdata, err := structpb.NewStruct(map[string]any{
		"category":         "Electronics",
		"customer_email":   "customer1@example.com",
		"id":               1,
		"product_id":       101,
		"product_name":     "Example Laptop 1",
		"product_type":     "Laptop",
		"shipping_address": "123 Main St, Cityville",
		"stock":            true,
	})
	require.NoError(t, err)

	return &opencdcv1.Record{
		Position:  []byte("position-1"),
		Operation: opencdcv1.Operation_OPERATION_CREATE,
		Metadata: map[string]string{
			"conduit.source.connector.id": "connector-1",
			"opencdc.readAt":              "1703019966257132000",
			"opencdc.version":             "v1",
		},
		Key: &opencdcv1.Data{
			Data: &opencdcv1.Data_StructuredData{
				StructuredData: keydata,
			},
		},
		Payload: &opencdcv1.Change{
			Before: nil,
			After: &opencdcv1.Data{
				Data: &opencdcv1.Data_StructuredData{
					StructuredData: afterdata,
				},
			},
		},
	}
}
