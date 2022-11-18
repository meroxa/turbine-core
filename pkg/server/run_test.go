package server

import (
	"context"
	"errors"
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

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
