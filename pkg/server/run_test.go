package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	pb "github.com/meroxa/turbine-core/v2/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/v2/pkg/app"
	"github.com/meroxa/turbine-core/v2/pkg/ir"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Init(t *testing.T) {
	ctx := context.Background()
	tempdir := t.TempDir()
	tests := []struct {
		desc    string
		setup   func() *pb.InitRequest
		wantErr error
	}{
		{
			desc:    "fails with invalid app name",
			wantErr: errors.New("invalid InitRequest.AppName: value length must be at least 1 runes"),
			setup: func() *pb.InitRequest {
				return &pb.InitRequest{
					ConfigFilePath: "/foo/bar",
					Language:       pb.Language_GOLANG,
				}
			},
		},
		{
			desc:    "fails with invalid config file path",
			wantErr: errors.New("invalid InitRequest.ConfigFilePath: value length must be at least 1 runes"),
			setup: func() *pb.InitRequest {
				return &pb.InitRequest{
					AppName:  "turbine-app",
					Language: pb.Language_GOLANG,
				}
			},
		},
		{
			desc:    "fails with invalid lang",
			wantErr: errors.New("invalid InitRequest.Language: value must be one of the defined enum values"),
			setup: func() *pb.InitRequest {
				return &pb.InitRequest{
					AppName:        "turbine-app",
					ConfigFilePath: "/foo/bar",
					Language:       101221,
				}
			},
		},
		{
			desc:    "fails to load app config",
			wantErr: errors.New("no such file or directory"),
			setup: func() *pb.InitRequest {
				return &pb.InitRequest{
					AppName:        "test-app",
					ConfigFilePath: "/nonexistingapp",
				}
			},
		},
		{
			desc: "success",
			setup: func() *pb.InitRequest {
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

				return &pb.InitRequest{
					AppName:        "app",
					ConfigFilePath: tempdir,
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			req := tc.setup()
			s := &runService{}

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
		setup   func() *pb.AddSourceRequest
		wantErr error
	}{
		{
			desc: "fails on invalid name",
			setup: func() *pb.AddSourceRequest {
				return &pb.AddSourceRequest{}
			},
			wantErr: errors.New("invalid AddSourceRequest.Name: value length must be at least 1 runes"),
		},
		{
			desc: "success",
			setup: func() *pb.AddSourceRequest {
				return &pb.AddSourceRequest{
					Name: "source-name",
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			s := &runService{}
			req := tc.setup()

			r, err := s.AddSource(ctx, req)
			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, r.StreamName, req.Name)
				}
			}
		})
	}
}

func Test_ReadRecords(t *testing.T) {
	ctx := context.Background()
	tempdir := t.TempDir()
	tests := []struct {
		desc        string
		srv         *runService
		setup       func() *pb.ReadRecordsRequest
		wantRecords *pb.ReadRecordsResponse
		wantErr     error
	}{
		{
			desc:    "fails when source is missing",
			srv:     &runService{},
			wantErr: errors.New("invalid ReadRecordsRequest.SourceStream: value length must be at least 1 runes"),
			setup: func() *pb.ReadRecordsRequest {
				return &pb.ReadRecordsRequest{}
			},
		},
		{
			desc: "fails on missing fixture file",
			srv: &runService{
				appPath: tempdir,
				config: app.Config{
					Fixtures: map[string]string{
						"resource": "fixture.json",
					},
				},
			},
			wantErr: errors.New("no such file or directory"),
			setup: func() *pb.ReadRecordsRequest {
				return &pb.ReadRecordsRequest{
					SourceStream: "resource",
				}
			},
		},
		{
			desc: "success",
			srv: &runService{
				appPath: path.Join(tempdir),
				config: app.Config{
					Fixtures: map[string]string{
						"source": "fixture.json",
					},
				},
			},
			wantRecords: &pb.ReadRecordsResponse{
				StreamRecords: &pb.StreamRecords{
					StreamName: "source",
					Records: []*pb.Record{
						{
							Key:   "1",
							Value: []byte(`{"message":"hello"}`),
						},
					},
				},
			},
			setup: func() *pb.ReadRecordsRequest {
				file := path.Join(tempdir, "fixture.json")
				require.NoError(
					t,
					os.WriteFile(
						file,
						[]byte(`[{
								"key": "1",
								"value": {"message":"hello"},
								"timestamp": "1662758822"
							}]`),
						0o644,
					),
				)
				return &pb.ReadRecordsRequest{SourceStream: "source"}
			},
		},
		{
			desc: "wrong fixture source name",
			srv: &runService{
				appPath: path.Join(tempdir),
				config: app.Config{
					Fixtures: map[string]string{
						"source123": "fixture.json",
					},
				},
			},
			wantErr: errors.New("no fixture file found for source pg"),
			setup: func() *pb.ReadRecordsRequest {
				file := path.Join(tempdir, "fixture.json")
				require.NoError(
					t,
					os.WriteFile(
						file,
						[]byte(`[{
								"key": "1",
								"value": {"message":"hello"},
								"timestamp": "1662758822"
							}]`),
						0o644,
					),
				)
				return &pb.ReadRecordsRequest{SourceStream: "pg"}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			req := tc.setup()

			c, err := tc.srv.ReadRecords(ctx, req)
			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, c.StreamRecords.StreamName, tc.wantRecords.StreamRecords.StreamName)
					assert.Equal(t, len(c.StreamRecords.Records), len(tc.wantRecords.StreamRecords.Records))
					assert.Equal(t, c.StreamRecords.Records[0].Key, tc.wantRecords.StreamRecords.Records[0].Key)
					assert.Equal(t, c.StreamRecords.Records[0].Value, tc.wantRecords.StreamRecords.Records[0].Value)
				}
			}
		})
	}
}

func Test_AddDestination(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		desc    string
		setup   func() *pb.AddDestinationRequest
		wantErr error
	}{
		{
			desc: "fails on invalid name",
			setup: func() *pb.AddDestinationRequest {
				return &pb.AddDestinationRequest{}
			},
			wantErr: errors.New("invalid AddDestinationRequest.Name: value length must be at least 1 runes"),
		},
		{
			desc: "success",
			setup: func() *pb.AddDestinationRequest {
				return &pb.AddDestinationRequest{
					Name: "destination-name",
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			s := &runService{}
			req := tc.setup()

			r, err := s.AddDestination(ctx, req)
			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, r.StreamName, req.Name)
				}
			}
		})
	}
}

func Test_WriteRecords(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		desc    string
		setup   func() *pb.WriteRecordsRequest
		wantErr error
	}{
		{
			desc:    "fails when destinationID is missing",
			wantErr: errors.New("invalid WriteRecordsRequest.DestinationID: value length must be at least 1 runes"),
			setup: func() *pb.WriteRecordsRequest {
				return &pb.WriteRecordsRequest{}
			},
		},
		{
			desc:    "fails when streamRecords is missing",
			wantErr: errors.New("invalid WriteRecordsRequest.StreamRecords: value is required"),
			setup: func() *pb.WriteRecordsRequest {
				return &pb.WriteRecordsRequest{
					DestinationID: "stream-destination",
				}
			},
		},
		{
			desc: "success",
			setup: func() *pb.WriteRecordsRequest {
				return &pb.WriteRecordsRequest{
					DestinationID: "destination-stream",
					StreamRecords: &pb.StreamRecords{
						StreamName: "source",
						Records: []*pb.Record{
							{
								Key:   "1",
								Value: []byte(`{"1":"record-value"}`),
							},
						},
					},
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			s := &runService{}
			req := tc.setup()

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
				assert.Contains(t, output, `{"1":"record-value"}`)
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
		setup   func() *pb.ProcessRecordsRequest
		wantErr error
	}{
		{
			desc: "fails on missing process",
			setup: func() *pb.ProcessRecordsRequest {
				return &pb.ProcessRecordsRequest{}
			},
			wantErr: errors.New("invalid ProcessRecordsRequest.Process: value is required"),
		},
		{
			desc: "fails on missing streamRecords",
			setup: func() *pb.ProcessRecordsRequest {
				return &pb.ProcessRecordsRequest{
					Process: &pb.ProcessRecordsRequest_Process{Name: "my-process"},
				}
			},
			wantErr: errors.New("invalid ProcessRecordsRequest.StreamRecords: value is required"),
		},
		{
			desc: "success",
			setup: func() *pb.ProcessRecordsRequest {
				return &pb.ProcessRecordsRequest{
					Process: &pb.ProcessRecordsRequest_Process{Name: "my-process"},
					StreamRecords: &pb.StreamRecords{
						StreamName: "my-stream",
					},
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			s := &runService{}
			req := tc.setup()

			c, err := s.ProcessRecords(ctx, req)
			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, c.StreamRecords.StreamName, req.StreamRecords.StreamName)
				}
			}
		})
	}
}
