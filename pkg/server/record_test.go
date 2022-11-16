package server

import (
	"context"
	"testing"

	"github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/pkg/ir"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	var (
		ctx  = context.Background()
		s    = NewRecordService()
		want = ir.DeploymentSpec{
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
		}
	)

	res, err := s.Init(ctx, &core.InitRequest{
		AppName:        "test-ruby",
		ConfigFilePath: "path/to/ruby",
		Language:       core.InitRequest_RUBY,
		GitSHA:         "gitsha",
		TurbineVersion: "0.1.0",
	})
	require.Nil(t, err)
	require.Equal(t, empty(), res)
	require.Equal(t, want, s.deploymentSpec)
}

func TestGetResource(t *testing.T) {
	var (
		ctx = context.Background()
		s   = NewRecordService()
	)

	res, err := s.GetResource(ctx, &core.GetResourceRequest{
		Name: "pg",
	})
	require.Nil(t, err)
	require.Equal(t, &core.Resource{Name: "pg"}, res)
	require.Equal(t, []*core.Resource{{Name: "pg"}}, s.resources)
}

func TestReadCollection(t *testing.T) {
	tests := []struct {
		description     string
		populateService func(*recordService) *recordService
		req             *core.ReadCollectionRequest
		want            ir.DeploymentSpec
		errMsg          string
	}{
		{
			description: "empty request",
			req:         &core.ReadCollectionRequest{},
			errMsg:      "please provide a collection name to 'read'",
		},
		{
			description: "recordService has existing source connector",
			req: &core.ReadCollectionRequest{
				Collection: "accounts",
				Resource: &core.Resource{
					Name: "pg",
				},
				Configs: nil,
			},
			populateService: func(s *recordService) *recordService {
				s.deploymentSpec.Connectors = []ir.ConnectorSpec{
					{
						Collection: "accounts",
						Resource:   "pg",
						Type:       ir.ConnectorSource,
					},
				}
				return s
			},
			errMsg: "only one call to 'read' is allowed per Meroxa Data Application",
		},
		{
			description: "successfully store source information",
			req: &core.ReadCollectionRequest{
				Collection: "accounts",
				Resource: &core.Resource{
					Name: "pg",
				},
				Configs: nil,
			},
			want: ir.DeploymentSpec{
				Connectors: []ir.ConnectorSpec{
					{
						Collection: "accounts",
						Resource:   "pg",
						Type:       ir.ConnectorSource,
						Config:     map[string]interface{}{},
					},
				},
			},
		},
		{
			description: "successfully store source information with config",
			req: &core.ReadCollectionRequest{
				Collection: "accounts",
				Resource: &core.Resource{
					Name: "pg",
				},
				Configs: &core.ResourceConfigs{
					ResourceConfig: []*core.ResourceConfig{
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
						Collection: "accounts",
						Resource:   "pg",
						Type:       ir.ConnectorSource,
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
			if test.populateService != nil {
				s = test.populateService(s)
			}

			res, err := s.ReadCollection(ctx, test.req)
			if test.errMsg != "" {
				require.EqualError(t, err, test.errMsg)
			} else {
				require.Nil(t, err)
				require.Equal(t, &core.Collection{}, res)
				require.Equal(t, test.want, s.deploymentSpec)
			}
		})
	}

}

func TestWriteCollectionToResource(t *testing.T) {
	tests := []struct {
		description     string
		populateService func(*recordService) *recordService
		req             *core.WriteCollectionRequest
		want            ir.DeploymentSpec
		errMsg          string
	}{
		{
			description: "empty request",
			req:         &core.WriteCollectionRequest{},
			errMsg:      "please provide a collection name to 'write'",
		},
		{
			description: "recordService has existing connector",
			req: &core.WriteCollectionRequest{
				TargetCollection: "accounts_copy",
				Resource: &core.Resource{
					Name: "pg",
				},
				Configs: nil,
			},
			populateService: func(s *recordService) *recordService {
				s.deploymentSpec.Connectors = []ir.ConnectorSpec{
					{
						Collection: "accounts",
						Resource:   "mongo",
						Type:       ir.ConnectorDestination,
					},
				}
				return s
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
			req: &core.WriteCollectionRequest{
				TargetCollection: "accounts_copy",
				Resource: &core.Resource{
					Name: "pg",
				},
				Configs: &core.ResourceConfigs{
					ResourceConfig: []*core.ResourceConfig{
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
			if test.populateService != nil {
				s = test.populateService(s)
			}

			res, err := s.WriteCollectionToResource(ctx, test.req)
			if test.errMsg != "" {
				require.EqualError(t, err, test.errMsg)
			} else {
				require.Nil(t, err)
				require.Equal(t, empty(), res)
				require.Equal(t, test.want, s.deploymentSpec)
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

	res, err := s.AddProcessToCollection(ctx,
		&core.ProcessCollectionRequest{
			Process: &core.Process{
				Name: "synchronize",
			},
		})
	require.Nil(t, err)
	require.Equal(t, &pb.Collection{}, res)
	require.Equal(t, want, s.deploymentSpec)
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
		&core.Secret{
			Name:  "api_key",
			Value: "secret_key",
		})
	require.Nil(t, err)
	require.Equal(t, empty(), res)

	res, err = s.RegisterSecret(ctx,
		&core.Secret{
			Name:  "another_key",
			Value: "key",
		})
	require.Nil(t, err)
	require.Equal(t, empty(), res)

	require.Equal(t, want, s.deploymentSpec)
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
					},
					{
						Name: "mongo",
					},
				}
				return s
			},
			want: &pb.ListResourcesResponse{
				Resources: []*pb.Resource{
					{
						Name: "pg",
					},
					{
						Name: "mongo",
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
