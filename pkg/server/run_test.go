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
	s := &runService{}
	r, err := s.GetResource(context.Background(), &pb.GetResourceRequest{
		Name: "my-resource",
	})

	if assert.NoError(t, err) {
		assert.Equal(t, r, &pb.Resource{
			Name: "my-resource",
		})
	}
}

func Test_AddProccessToCollection(t *testing.T) {
	s := &runService{}
	c, err := s.AddProcessToCollection(context.Background(), &pb.ProcessCollectionRequest{
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
	})

	if assert.NoError(t, err) {
		assert.Equal(t, c, &pb.Collection{
			Name:   "my-collection",
			Stream: "my-stream",
			Records: []*pb.Record{
				{
					Key:   "key1",
					Value: []byte("val1"),
				},
			},
		})
	}
}

func TestRunService_RegisterSecret(t *testing.T) {
	tests := []struct {
		description string
		envVars     map[string]string
		secretWant  *pb.Secret
		wantErr     error
	}{
		{
			description: "when secret exists",
			envVars: map[string]string{
				"TEST_ENV_VAR": "my-value",
			},
			secretWant: &pb.Secret{
				Name: "TEST_ENV_VAR",
			},
			wantErr: nil,
		},
		{
			description: "when secret doesn't exist",
			envVars: map[string]string{
				"TEST_ENV_VAR": "my-value",
			},
			secretWant: &pb.Secret{
				Name: "TEST_FOO",
			},
			wantErr: errors.New("secret is invalid or not set"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			var (
				ctx = context.Background()
				s   = NewRunService()
			)

			cleanup := envSet(tc.envVars)

			_, gotErr := s.RegisterSecret(ctx, tc.secretWant)
			if tc.wantErr != nil {
				assert.Equal(t, gotErr, tc.wantErr)
			} else {
				assert.NoError(t, gotErr)
			}

			t.Cleanup(cleanup)
		})
	}
}

// envSet will take care of setting and unsetting environment variables during cleanup
func envSet(envs map[string]string) (cleanup func()) {
	for name, value := range envs {
		_ = os.Setenv(name, value)
	}

	return func() {
		for name := range envs {
			_ = os.Unsetenv(name)
		}
	}
}
