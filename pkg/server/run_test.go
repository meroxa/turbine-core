package server

import (
	"context"
	"errors"
	"os"
	"path"
	"testing"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/app"

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
						[]byte(`{
							"name": "app",
							"language": "ruby",
							"environment": "common",
							"resources": {
								"demopg": "fixtures/demo.json"
							}
						}`),
						0644,
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
					Name:        "app",
					Environment: "common",
					Pipeline:    "turbine-pipeline-app",
					Resources: map[string]string{
						"demopg": "fixtures/demo.json",
					},
				})
			}
		})
	}
}

func Test_GetResource(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		desc    string
		setup   func() *pb.GetResourceRequest
		wantErr error
	}{
		{
			desc: "fails on invalid name",
			setup: func() *pb.GetResourceRequest {
				return &pb.GetResourceRequest{}
			},
			wantErr: errors.New("invalid GetResourceRequest.Name: value length must be at least 1 runes"),
		},
		{
			desc: "success",
			setup: func() *pb.GetResourceRequest {
				return &pb.GetResourceRequest{
					Name: "my-resource",
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			s := &runService{}
			req := tc.setup()

			r, err := s.GetResource(ctx, req)
			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, r.Name, req.Name)
				}
			}
		})
	}
}

func Test_AddProccessToCollection(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		desc    string
		setup   func() *pb.ProcessCollectionRequest
		wantErr error
	}{
		{
			desc: "fails on missing process",
			setup: func() *pb.ProcessCollectionRequest {
				return &pb.ProcessCollectionRequest{}
			},
			wantErr: errors.New("invalid ProcessCollectionRequest.Process: value is required"),
		},
		{
			desc: "fails on missing collection",
			setup: func() *pb.ProcessCollectionRequest {
				return &pb.ProcessCollectionRequest{
					Process: &pb.ProcessCollectionRequest_Process{},
				}
			},
			wantErr: errors.New("invalid ProcessCollectionRequest.Collection: value is required"),
		},
		{
			desc: "success",
			setup: func() *pb.ProcessCollectionRequest {
				return &pb.ProcessCollectionRequest{
					Process: &pb.ProcessCollectionRequest_Process{
						Name: "my-process",
					},
					Collection: &pb.Collection{
						Name:   "my-collection",
						Stream: "my-stream",
						Records: []*pb.Record{
							{
								Key:   "key1",
								Value: []byte("val1"),
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

			c, err := s.AddProcessToCollection(ctx, req)
			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, c, req.Collection)
				}
			}
		})
	}
}

func Test_RegisterSecret(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		desc    string
		setup   func() *pb.Secret
		wantErr error
	}{
		{
			desc:    "fails when secret name is invalid",
			wantErr: errors.New("invalid Secret.Name: value length must be at least 1 runes"),
			setup: func() *pb.Secret {
				return &pb.Secret{
					Value: "secret-name",
				}
			},
		},
		{
			desc:    "fails when secret value is invalid",
			wantErr: errors.New("invalid Secret.Value: value length must be at least 1 runes"),
			setup: func() *pb.Secret {
				return &pb.Secret{
					Name: "secret-value",
				}
			},
		},
		{
			desc: "success",
			setup: func() *pb.Secret {
				return &pb.Secret{
					Name:  "secret-name",
					Value: "secret-value",
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			s := &runService{}
			req := tc.setup()

			_, err := s.RegisterSecret(ctx, req)
			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_ReadCollection(t *testing.T) {
	ctx := context.Background()
	tempdir := t.TempDir()
	tests := []struct {
		desc           string
		srv            *runService
		setup          func() *pb.ReadCollectionRequest
		wantCollection *pb.Collection
		wantErr        error
	}{
		{
			desc:    "fails when resource is missing",
			srv:     &runService{},
			wantErr: errors.New("invalid ReadCollectionRequest.Resource: value is required"),
			setup: func() *pb.ReadCollectionRequest {
				return &pb.ReadCollectionRequest{
					Collection: "resource-collection",
				}
			},
		},
		{
			desc:    "fails when collection name is missing",
			srv:     &runService{},
			wantErr: errors.New("invalid ReadCollectionRequest.Collection: value length must be at least 1 runes"),
			setup: func() *pb.ReadCollectionRequest {
				return &pb.ReadCollectionRequest{
					Resource: &pb.Resource{
						Name: "resource",
					},
				}
			},
		},
		{
			desc: "fails on missing fixture file",
			srv: &runService{
				appPath: tempdir,
				config: app.Config{
					Resources: map[string]string{
						"resource": "fixture.json",
					},
				},
			},
			wantErr: errors.New("no such file or directory"),
			setup: func() *pb.ReadCollectionRequest {
				return &pb.ReadCollectionRequest{
					Resource: &pb.Resource{
						Name: "resource",
					},
					Collection: "resource-collection",
				}
			},
		},
		{
			desc: "success",
			srv: &runService{
				appPath: path.Join(tempdir),
				config: app.Config{
					Resources: map[string]string{
						"resource": "fixture.json",
					},
				},
			},
			wantCollection: &pb.Collection{
				Name: "events",
				Records: []*pb.Record{
					{
						Key:   "1",
						Value: []byte(`{"message":"hello"}`),
					},
				},
			},
			setup: func() *pb.ReadCollectionRequest {
				file := path.Join(tempdir, "fixture.json")
				require.NoError(
					t,
					os.WriteFile(
						file,
						[]byte(`{
							"events": [{
								"key": "1",
								"value": {"message":"hello"},
								"timestamp": "1662758822"
							}]
						}`),
						0644,
					),
				)
				return &pb.ReadCollectionRequest{
					Resource: &pb.Resource{
						Name: "resource",
					},
					Collection: "events",
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			req := tc.setup()

			c, err := tc.srv.ReadCollection(ctx, req)
			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, c.Name, tc.wantCollection.Name)
					assert.Equal(t, len(c.Records), len(tc.wantCollection.Records))
					assert.Equal(t, c.Records[0].Key, tc.wantCollection.Records[0].Key)
					assert.Equal(t, c.Records[0].Value, tc.wantCollection.Records[0].Value)
				}
			}
		})
	}
}

func Test_WriteCollectionToResource(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		desc    string
		setup   func() *pb.WriteCollectionRequest
		wantErr error
	}{}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			s := &runService{}
			req := tc.setup()

			_, err := s.WriteCollectionToResource(ctx, req)
			if tc.wantErr != nil {
				assert.ErrorContains(t, err, tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
