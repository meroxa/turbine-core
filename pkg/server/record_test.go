package server

import (
	"context"
	"errors"
	"fmt"
	"testing"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/ir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	testCases := []struct {
		test    string
		spec    *ir.DeploymentSpec
		request *pb.InitRequest
		want    error
	}{
		{
			test: "Init successful with correct language",
			spec: &ir.DeploymentSpec{
				Definition: ir.DefinitionSpec{
					GitSha: "gitsha",
					Metadata: ir.MetadataSpec{
						Turbine: ir.TurbineSpec{
							Language: ir.Ruby,
							Version:  "0.1.0",
						},
						SpecVersion: ir.LatestSpecVersion,
					},
				},
			},
			request: &pb.InitRequest{
				AppName:        "test-ruby",
				ConfigFilePath: "path/to/ruby",
				Language:       pb.Language_RUBY,
				GitSHA:         "gitsha",
				TurbineVersion: "0.1.0",
			},
			want: nil,
		},
		{
			test: "Init error with incorrect language",
			spec: &ir.DeploymentSpec{
				Definition: ir.DefinitionSpec{
					GitSha: "gitsha",
					Metadata: ir.MetadataSpec{
						Turbine: ir.TurbineSpec{
							Language: "emoji",
							Version:  "0.1.0",
						},
						SpecVersion: ir.LatestSpecVersion,
					},
				},
			},
			request: &pb.InitRequest{
				AppName:        "test-emoji",
				ConfigFilePath: "path/to/emoji",
				Language:       101221,
				GitSHA:         "gitsha",
				TurbineVersion: "0.1.0",
			},
			want: errors.New("invalid InitRequest.Language: value must be one of the defined enum values"),
		},
	}

	for _, test := range testCases {
		t.Run(test.test, func(t *testing.T) {
			var (
				ctx = context.Background()
				s   = NewSpecBuilderService()
			)
			res, err := s.Init(ctx, test.request)

			if test.want == nil {
				require.Nil(t, err)
				require.Equal(t, empty(), res)
				require.Equal(t, test.spec.Functions, s.spec.Functions)
				require.Equal(t, test.spec.Connectors, s.spec.Connectors)
				require.Equal(t, test.spec.Streams, s.spec.Streams)
			} else {
				require.ErrorContains(t, err, test.want.Error())
			}

		})
	}

}

func TestGetResource(t *testing.T) {
	var (
		ctx = context.Background()
		s   = NewSpecBuilderService()
	)

	res, err := s.GetResource(ctx, &pb.GetResourceRequest{
		Name: "pg",
	})
	require.Nil(t, err)
	require.Equal(t, &pb.Resource{Name: "pg"}, res)
}

func TestReadCollection(t *testing.T) {
	tests := []struct {
		description     string
		populateService func(*specBuilderService) *specBuilderService
		req             *pb.ReadCollectionRequest
		want            *ir.DeploymentSpec
		errMsg          string
	}{

		{
			description: "successfully store source information",
			req: &pb.ReadCollectionRequest{
				Collection: "accounts",
				Resource: &pb.Resource{
					Name:       "pg",
					Source:     true,
					Collection: "accounts",
				},
				Configs: nil,
			},
		},
		{
			description: "successfully store source information with config",
			req: &pb.ReadCollectionRequest{
				Collection: "accounts",
				Resource: &pb.Resource{
					Name:       "pg",
					Source:     true,
					Collection: "accounts",
				},
				Configs: &pb.Configs{
					Config: []*pb.Config{
						{
							Field: "config",
							Value: "value",
						},
						{
							Field: "another_config",
							Value: "another_value",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			var (
				ctx = context.Background()
				s   = NewSpecBuilderService()
			)
			if test.populateService != nil {
				s = test.populateService(s)
			}

			res, err := s.ReadCollection(ctx, test.req)
			if test.errMsg != "" {
				require.EqualError(t, err, test.errMsg)
			} else {
				require.Nil(t, err)
				require.Equal(t, test.req.Resource.Collection, res.Name)
				require.NotEmpty(t, s.spec.Connectors)
				require.Equal(t, s.spec.Connectors[0].Collection, res.Name)
				require.Equal(t, s.spec.Connectors[0].UUID, res.Stream)
				require.Equal(t, s.spec.Connectors[0].Type, ir.ConnectorType("source"))
			}
		})
	}

}

func TestWriteCollectionToResource(t *testing.T) {
	tests := []struct {
		description     string
		populateService func(*specBuilderService) *specBuilderService
		req             *pb.WriteCollectionRequest
		want            *ir.DeploymentSpec
		errMsg          string
	}{
		{
			description: "empty request",
			req:         &pb.WriteCollectionRequest{},
			errMsg:      "invalid WriteCollectionRequest.Resource: value is required",
		},
		{
			description: "specBuilderService has existing connector",
			req: &pb.WriteCollectionRequest{
				Resource: &pb.Resource{
					Name: "pg",
				},
				TargetCollection: "accounts_copy",
				Configs:          nil,
			},
			want: &ir.DeploymentSpec{
				Connectors: []ir.ConnectorSpec{
					{
						Collection: "accounts",
						Resource:   "mongo",
						Type:       ir.ConnectorDestination,
					},
					{
						Collection: "accounts_copy",
						Resource:   "pg",
						Type:       ir.ConnectorDestination,
						Config:     map[string]interface{}{},
					},
				},
			},
		},
		{
			description: "successfully store destination information with config",
			req: &pb.WriteCollectionRequest{
				Resource: &pb.Resource{
					Name: "pg",
				},
				TargetCollection: "accounts_copy",
				Configs: &pb.Configs{
					Config: []*pb.Config{
						{
							Field: "config",
							Value: "value",
						},
						{
							Field: "another_config",
							Value: "another_value",
						},
					},
				},
			},
			want: &ir.DeploymentSpec{
				Connectors: []ir.ConnectorSpec{
					{
						Collection: "accounts_copy",
						Resource:   "pg",
						Type:       ir.ConnectorDestination,
						Config: map[string]interface{}{
							"config":         "value",
							"another_config": "another_value",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			var (
				ctx = context.Background()
				s   = NewSpecBuilderService()
			)

			source, err := s.ReadCollection(ctx,
				&pb.ReadCollectionRequest{
					Collection: "pg_2",
					Resource: &pb.Resource{
						Name:       "pg_2",
						Source:     true,
						Collection: "pg_2",
					},
					Configs: nil,
				},
			)
			assert.NoError(t, err)

			test.req.SourceCollection = source

			res, err := s.WriteCollectionToResource(ctx, test.req)
			if test.errMsg != "" {
				require.EqualError(t, err, test.errMsg)
			} else {

				require.Nil(t, err)
				require.Equal(t, empty(), res)
				require.NotEmpty(t, s.spec.Streams)
				require.NotEmpty(t, s.spec.Connectors)
				require.Equal(t, s.spec.Streams[0].FromUUID, source.Stream)
				require.Equal(t, s.spec.Streams[0].ToUUID, s.spec.Connectors[1].UUID)
			}
		})
	}

}

func TestAddProcessToCollection(t *testing.T) {
	var (
		ctx  = context.Background()
		s    = NewSpecBuilderService()
		want = &ir.DeploymentSpec{
			Functions: []ir.FunctionSpec{
				{
					Name: "synchronize",
				},
			},
		}
	)

	read, err := s.ReadCollection(ctx,
		&pb.ReadCollectionRequest{

			Collection: "accounts",
			Resource: &pb.Resource{
				Name:       "pg",
				Source:     true,
				Collection: "accounts",
			},
			Configs: nil,
		},
	)
	assert.NoError(t, err)

	res, err := s.AddProcessToCollection(ctx,
		&pb.ProcessCollectionRequest{
			Process: &pb.ProcessCollectionRequest_Process{
				Name: "synchronize",
			},
			Collection: read,
		},
	)

	require.Nil(t, err)
	require.NotEmpty(t, res)
	require.NotEmpty(t, s.spec.Functions)
	require.Equal(t, s.spec.Functions[0].Name, want.Functions[0].Name)
	require.Equal(t, s.spec.Streams[0].FromUUID, read.Stream)
	require.Equal(t, s.spec.Streams[0].ToUUID, res.Stream)

}

func TestRegisterSecret(t *testing.T) {
	var (
		ctx  = context.Background()
		s    = NewSpecBuilderService()
		want = &ir.DeploymentSpec{
			Secrets: map[string]string{
				"api_key":     "secret_key",
				"another_key": "key",
			},
		}
	)

	res, err := s.RegisterSecret(ctx,
		&pb.Secret{
			Name:  "api_key",
			Value: "secret_key",
		})
	require.Nil(t, err)
	require.Equal(t, empty(), res)

	res, err = s.RegisterSecret(ctx,
		&pb.Secret{
			Name:  "another_key",
			Value: "key",
		})
	require.Nil(t, err)
	require.Equal(t, empty(), res)
	require.Equal(t, want.Secrets, s.spec.Secrets)

}

func TestHasFunctions(t *testing.T) {
	tests := []struct {
		description     string
		populateService func(*specBuilderService) *specBuilderService
		want            bool
	}{
		{
			description: "service with no functions",
			want:        false,
		},
		{
			description: "service with function",
			populateService: func(s *specBuilderService) *specBuilderService {
				s.spec.Functions = []ir.FunctionSpec{
					{
						Name: "addition",
					},
				}
				return s
			},
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			var (
				ctx = context.Background()
				s   = NewSpecBuilderService()
			)
			if test.populateService != nil {
				s = test.populateService(s)
			}

			res, err := s.HasFunctions(ctx, empty())
			require.Nil(t, err)
			require.Equal(t, test.want, res.Value)
		})
	}
}

func TestListResources(t *testing.T) {
	tests := []struct {
		description     string
		populateService func(*specBuilderService) *specBuilderService
		want            *pb.ListResourcesResponse
	}{
		{
			description: "service with no resources",
			want:        &pb.ListResourcesResponse{},
		},
		{
			description: "service with resources",
			populateService: func(s *specBuilderService) *specBuilderService {
				s.resources = []*pb.Resource{
					{
						Name: "pg",

						Source:     true,
						Collection: "in",
					},
					{
						Name: "mongo",

						Destination: true,
						Collection:  "out",
					},
				}
				return s
			},
			want: &pb.ListResourcesResponse{
				Resources: []*pb.Resource{
					{
						Name: "pg",

						Source:     true,
						Collection: "in",
					},
					{
						Name: "mongo",

						Destination: true,
						Collection:  "out",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			var (
				ctx = context.Background()
				s   = NewSpecBuilderService()
			)
			if test.populateService != nil {
				s = test.populateService(s)
			}

			res, err := s.ListResources(ctx, empty())
			require.Nil(t, err)
			require.Equal(t, test.want, res)
		})
	}
}

func TestGetSpec(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		description     string
		populateService func(*specBuilderService) *specBuilderService
		request         *pb.GetSpecRequest
		want            *ir.DeploymentSpec
		wantErr         error
	}{
		{
			description: "get spec with no function",
			populateService: func(s *specBuilderService) *specBuilderService {
				s.spec = exampleDeploymentSpec()
				s.spec.Streams = append(s.spec.Streams, ir.StreamSpec{
					UUID:     "1_3",
					FromUUID: "1",
					ToUUID:   "3",
					Name:     "1_3",
				})

				return s
			},
			request: &pb.GetSpecRequest{
				Image: "",
			},
			want: func() *ir.DeploymentSpec {
				s := exampleDeploymentSpec()
				s.Streams = append(s.Streams, ir.StreamSpec{
					UUID:     "1_3",
					FromUUID: "1",
					ToUUID:   "3",
					Name:     "1_3",
				})
				return s
			}(),
		},
		{
			description: "get spec with no function, set image",
			populateService: func(s *specBuilderService) *specBuilderService {
				s.spec = exampleDeploymentSpec()
				s.spec.Streams = append(s.spec.Streams, ir.StreamSpec{
					UUID:     "1_3",
					FromUUID: "1",
					ToUUID:   "3",
					Name:     "1_3",
				})
				return s
			},
			request: &pb.GetSpecRequest{
				Image: "some/image",
			},
			wantErr: fmt.Errorf("cannot set image without defined functions"),
		},
		{
			description: "get spec with function",
			populateService: func(s *specBuilderService) *specBuilderService {
				s.spec = exampleDeploymentSpec()
				s.spec.Functions = append(s.spec.Functions, ir.FunctionSpec{
					UUID:  "2",
					Name:  "function",
					Image: "some/image",
				})
				s.spec.Streams = append(s.spec.Streams, ir.StreamSpec{
					UUID:     "1_2",
					FromUUID: "1",
					ToUUID:   "2",
					Name:     "1_2",
				})
				s.spec.Streams = append(s.spec.Streams, ir.StreamSpec{
					UUID:     "2_3",
					FromUUID: "2",
					ToUUID:   "3",
					Name:     "2_3",
				})
				return s
			},
			request: &pb.GetSpecRequest{
				Image: "some/image",
			},
			want: func() *ir.DeploymentSpec {
				s := exampleDeploymentSpec()
				s.AddFunction(
					&ir.FunctionSpec{
						UUID:  "2",
						Name:  "function",
						Image: "some/image",
					},
				)
				s.Streams = append(s.Streams, ir.StreamSpec{
					UUID:     "1_2",
					FromUUID: "1",
					ToUUID:   "2",
					Name:     "1_2",
				})
				s.Streams = append(s.Streams, ir.StreamSpec{
					UUID:     "2_3",
					FromUUID: "2",
					ToUUID:   "3",
					Name:     "2_3",
				})

				return s
			}(),
			wantErr: nil,
		},
		{
			description: "get spec with function, overwrite image",
			populateService: func(s *specBuilderService) *specBuilderService {
				s.spec = exampleDeploymentSpec()
				s.spec.Functions = append(s.spec.Functions, ir.FunctionSpec{
					UUID:  "2",
					Name:  "function",
					Image: "some/image",
				})
				s.spec.Streams = append(s.spec.Streams, ir.StreamSpec{
					UUID:     "1_2",
					FromUUID: "1",
					ToUUID:   "2",
					Name:     "1_2",
				})
				s.spec.Streams = append(s.spec.Streams, ir.StreamSpec{
					UUID:     "2_3",
					FromUUID: "2",
					ToUUID:   "3",
					Name:     "2_3",
				})
				return s
			},
			request: &pb.GetSpecRequest{
				Image: "some/image",
			},
			want: func() *ir.DeploymentSpec {
				s := exampleDeploymentSpec()
				s.AddFunction(
					&ir.FunctionSpec{
						UUID:  "2",
						Name:  "function",
						Image: "some/image",
					},
				)

				s.Streams = append(s.Streams, ir.StreamSpec{
					UUID:     "1_2",
					FromUUID: "1",
					ToUUID:   "2",
					Name:     "1_2",
				})
				s.Streams = append(s.Streams, ir.StreamSpec{
					UUID:     "2_3",
					FromUUID: "2",
					ToUUID:   "3",
					Name:     "2_3",
				})

				return s
			}(),
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			s := test.populateService(NewSpecBuilderService())
			res, err := s.GetSpec(ctx, test.request)
			if test.wantErr == nil && err == nil {
				got, err := ir.Unmarshal(res.Spec)
				require.NoError(t, err)
				require.Equal(t, test.want.Connectors, got.Connectors)
				require.Equal(t, test.want.Functions, got.Functions)
				require.Equal(t, test.want.Secrets, got.Secrets)
				require.Equal(t, test.want.Streams, got.Streams)
			} else {
				require.Error(t, err)
				require.Equal(t, test.wantErr, err)
			}
		})
	}
}

func exampleDeploymentSpec() *ir.DeploymentSpec {
	return &ir.DeploymentSpec{
		Secrets: map[string]string{
			"a secret": "with value",
		},
		Connectors: []ir.ConnectorSpec{
			{
				UUID:       "1",
				Collection: "accounts",
				Resource:   "mongo",
				Type:       ir.ConnectorSource,
			},
			{
				UUID:       "3",
				Collection: "accounts_copy",
				Resource:   "pg",
				Type:       ir.ConnectorDestination,
				Config: map[string]interface{}{
					"config": "value",
				},
			},
		},
		Definition: ir.DefinitionSpec{
			GitSha: "gitsh",
			Metadata: ir.MetadataSpec{
				SpecVersion: "0.2.0",
				Turbine: ir.TurbineSpec{
					Language: ir.GoLang,
					Version:  "10",
				},
			},
		},
	}
}
