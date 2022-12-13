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
		spec    ir.DeploymentSpec
		request pb.InitRequest
		want    error
	}{
		{
			test: "Init successful with correct language",
			spec: ir.DeploymentSpec{
				Definition: ir.DefinitionSpec{
					GitSha: "gitsha",
					Metadata: ir.MetadataSpec{
						Turbine: ir.TurbineSpec{
							Language: "ruby",
							Version:  "0.1.0",
						},
						SpecVersion: ir.LatestSpecVersion,
					},
				},
			},
			request: pb.InitRequest{
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
			spec: ir.DeploymentSpec{
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
			request: pb.InitRequest{
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
				s   = NewRecordService()
			)
			res, err := s.Init(ctx, &test.request)

			if test.want == nil {
				require.Nil(t, err)
				require.Equal(t, empty(), res)
				require.Equal(t, test.spec.Functions, s.deploymentSpec.Functions)
				require.Equal(t, test.spec.Connectors, s.deploymentSpec.Connectors)
				require.Equal(t, test.spec.Streams, s.deploymentSpec.Streams)
			} else {
				require.ErrorContains(t, err, test.want.Error())
			}

		})
	}

}

func TestGetResource(t *testing.T) {
	var (
		ctx = context.Background()
		s   = NewRecordService()
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
		populateService func(*recordService) *recordService
		req             *pb.ReadCollectionRequest
		want            ir.DeploymentSpec
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
				s   = NewRecordService()
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
				require.NotEmpty(t, s.deploymentSpec.Connectors)
				require.Equal(t, s.deploymentSpec.Connectors[0].Collection, res.Name)
				require.Equal(t, s.deploymentSpec.Connectors[0].UUID, res.Stream)
				require.Equal(t, s.deploymentSpec.Connectors[0].Type, ir.ConnectorType("source"))
			}
		})
	}

}

func TestWriteCollectionToResource(t *testing.T) {
	tests := []struct {
		description     string
		populateService func(*recordService) *recordService
		req             *pb.WriteCollectionRequest
		want            ir.DeploymentSpec
		errMsg          string
	}{
		{
			description: "empty request",
			req:         &pb.WriteCollectionRequest{},
			errMsg:      "please provide a collection name to 'write'",
		},
		{
			description: "recordService has existing connector",
			req: &pb.WriteCollectionRequest{
				TargetCollection: "accounts_copy",
				Configs:          nil,
			},
			want: ir.DeploymentSpec{
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
			want: ir.DeploymentSpec{
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
				s   = NewRecordService()
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
				require.NotEmpty(t, s.deploymentSpec.Streams)
				require.NotEmpty(t, s.deploymentSpec.Connectors)
				require.Equal(t, s.deploymentSpec.Streams[0].FromUUID, source.Stream)
				require.Equal(t, s.deploymentSpec.Streams[0].ToUUID, s.deploymentSpec.Connectors[1].UUID)
			}
		})
	}

}

func TestAddProcessToCollection(t *testing.T) {
	var (
		ctx  = context.Background()
		s    = NewRecordService()
		want = ir.DeploymentSpec{
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
	require.NotEmpty(t, s.deploymentSpec.Functions)
	require.Equal(t, s.deploymentSpec.Functions[0].Name, want.Functions[0].Name)
	require.Equal(t, s.deploymentSpec.Streams[0].FromUUID, read.Stream)
	require.Equal(t, s.deploymentSpec.Streams[0].ToUUID, res.Stream)

}

func TestRegisterSecret(t *testing.T) {
	var (
		ctx  = context.Background()
		s    = NewRecordService()
		want = ir.DeploymentSpec{
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
	require.Equal(t, want.Secrets, s.deploymentSpec.Secrets)

}

func TestHasFunctions(t *testing.T) {
	tests := []struct {
		description     string
		populateService func(*recordService) *recordService
		want            bool
	}{
		{
			description: "service with no functions",
			want:        false,
		},
		{
			description: "service with function",
			populateService: func(s *recordService) *recordService {
				s.deploymentSpec.Functions = []ir.FunctionSpec{
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
				s   = NewRecordService()
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
		populateService func(*recordService) *recordService
		want            *pb.ListResourcesResponse
	}{
		{
			description: "service with no resources",
			want:        &pb.ListResourcesResponse{},
		},
		{
			description: "service with resources",
			populateService: func(s *recordService) *recordService {
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
				s   = NewRecordService()
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
	tests := []struct {
		description     string
		populateService func(*recordService) *recordService
		request         *pb.GetSpecRequest
		want            ir.DeploymentSpec
		wantErr         error
	}{
		{
			description: "get spec with no function",
			populateService: func(s *recordService) *recordService {
				s.deploymentSpec = exampleDeploymentSpec()
				return s
			},
			want: exampleDeploymentSpec(),
		},
		{
			description: "get spec with no function, set image",
			populateService: func(s *recordService) *recordService {
				s.deploymentSpec = exampleDeploymentSpec()
				return s
			},
			request: &pb.GetSpecRequest{
				Image: "some/image",
			},
			wantErr: fmt.Errorf("cannot set function image since spec has no functions"),
		},
		{
			description: "get spec with function",
			populateService: func(s *recordService) *recordService {
				s.deploymentSpec = exampleDeploymentSpec()
				s.deploymentSpec.Functions = []ir.FunctionSpec{
					{
						Name: "function",
					},
				}
				return s
			},
			request: &pb.GetSpecRequest{
				Image: "some/image",
			},
			want: func() ir.DeploymentSpec {
				s := exampleDeploymentSpec()
				s.Functions = []ir.FunctionSpec{
					{
						Name:  "function",
						Image: "some/image",
					},
				}
				return s
			}(),
		},
		{
			description: "get spec with function, overwrite image",
			populateService: func(s *recordService) *recordService {
				s.deploymentSpec = exampleDeploymentSpec()
				s.deploymentSpec.Functions = []ir.FunctionSpec{
					{
						Name:  "function",
						Image: "existing/image",
					},
				}
				return s
			},
			request: &pb.GetSpecRequest{
				Image: "some/image",
			},
			want: func() ir.DeploymentSpec {
				s := exampleDeploymentSpec()
				s.Functions = []ir.FunctionSpec{
					{
						Name:  "function",
						Image: "some/image",
					},
				}
				return s
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			var (
				ctx = context.Background()
				s   = test.populateService(NewRecordService())
			)

			res, err := s.GetSpec(ctx, test.request)
			require.Equal(t, test.wantErr, err)
			if test.wantErr == nil {
				got, err := ir.Unmarshal(res.Spec)
				require.Nil(t, err)
				require.Equal(t, got, &test.want)
			}
		})
	}
}

func exampleDeploymentSpec() ir.DeploymentSpec {
	return ir.DeploymentSpec{
		Secrets: map[string]string{
			"a secret": "with value",
		},
		Connectors: []ir.ConnectorSpec{
			{
				Collection: "accounts",
				Resource:   "mongo",
				Type:       ir.ConnectorSource,
			},
			{
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
