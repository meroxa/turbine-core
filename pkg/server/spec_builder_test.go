package server

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/conduitio/conduit-commons/proto/opencdc/v1"
	"github.com/google/uuid"
	"github.com/meroxa/turbine-core/v2/pkg/ir"
	"github.com/meroxa/turbine-core/v2/proto/turbine/v2"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	testCases := []struct {
		test    string
		spec    *ir.DeploymentSpec
		request *turbinev2.InitRequest
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
			request: &turbinev2.InitRequest{
				AppName:        "test-ruby",
				ConfigFilePath: "path/to/ruby",
				Language:       turbinev2.Language_RUBY,
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
			request: &turbinev2.InitRequest{
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

func TestAddSource(t *testing.T) {
	tests := []struct {
		description     string
		populateService func(*SpecBuilderService) *SpecBuilderService
		req             *turbinev2.AddSourceRequest
		want            *ir.DeploymentSpec
		errMsg          string
	}{
		{
			description: "successfully store source information",
			req: &turbinev2.AddSourceRequest{
				Name: "my-source",
				Plugin: &turbinev2.Plugin{
					Name: "builtin:postgres@1.0.0",
				},
			},
		},
		{
			description: "successfully store source information with config",
			req: &turbinev2.AddSourceRequest{
				Name: "my-source",
				Plugin: &turbinev2.Plugin{
					Name: "builtin:postgres@1.0.0",
					Config: map[string]string{
						"collection":     "accounts",
						"config":         "value",
						"another_config": "another_value",
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

			res, err := s.AddSource(ctx, test.req)
			if test.errMsg != "" {
				require.EqualError(t, err, test.errMsg)
			} else {
				require.Nil(t, err)
				require.NotEmpty(t, s.spec.Connectors)
				require.Equal(t, s.spec.Connectors[0].Name, test.req.Name)
				require.Equal(t, s.spec.Connectors[0].UUID, res.StreamName)
				require.Equal(t, s.spec.Connectors[0].PluginType, ir.PluginSource)
			}
		})
	}
}

func TestReadRecords(t *testing.T) {
	var (
		ctx  = context.Background()
		s    = NewSpecBuilderService()
		uuid = uuid.New().String()
	)

	res, err := s.ReadRecords(ctx, &turbinev2.ReadRecordsRequest{
		SourceStream: uuid,
	})
	require.Nil(t, err)
	require.Equal(t, &turbinev2.ReadRecordsResponse{
		StreamRecords: &turbinev2.StreamRecords{
			StreamName: uuid,
		},
	}, res)
}

func TestAddDestination(t *testing.T) {
	tests := []struct {
		description     string
		populateService func(*SpecBuilderService) *SpecBuilderService
		req             *turbinev2.AddDestinationRequest
		want            *ir.DeploymentSpec
		errMsg          string
	}{
		{
			description: "empty request",
			req:         &turbinev2.AddDestinationRequest{},
			errMsg:      "invalid AddDestinationRequest.Name: value length must be at least 1 runes",
		},
		{
			description: "successfully store destination information with config",
			req: &turbinev2.AddDestinationRequest{
				Name: "my-destination",
				Plugin: &turbinev2.Plugin{
					Name: "builtin:postgres@1.0.0",
					Config: map[string]string{
						"collection": "accounts_copy",
					},
				},
			},
			want: &ir.DeploymentSpec{
				Connectors: []ir.ConnectorSpec{
					{
						Name:       "my-destination",
						PluginName: "builtin:postgres@1.0.0",
						PluginType: ir.PluginDestination,
						PluginConfig: map[string]string{
							"collection": "accounts_copy",
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

			_, err := s.AddDestination(ctx, test.req)

			if test.errMsg != "" {
				require.EqualError(t, err, test.errMsg)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, s.spec.Connectors)
				require.Equal(t, s.spec.Connectors[0].Name, test.req.Name)
				require.Equal(t, s.spec.Connectors[0].PluginType, ir.PluginDestination)
			}
		})
	}
}

func TestWriteRecords(t *testing.T) {
	tests := []struct {
		description     string
		populateService func(*SpecBuilderService) *SpecBuilderService
		req             *turbinev2.WriteRecordsRequest
		want            *ir.DeploymentSpec
		errMsg          string
	}{
		{
			description: "empty request",
			req:         &turbinev2.WriteRecordsRequest{},
			errMsg:      "invalid WriteRecordsRequest.DestinationID: value length must be at least 1 runes",
		},
		{
			description: "successfully store stream information",
			req: &turbinev2.WriteRecordsRequest{
				StreamRecords: &turbinev2.StreamRecords{
					Records: []*opencdcv1.Record{},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			var (
				ctx = context.Background()
				s   = NewSpecBuilderService()
				err error
			)

			asr, err := s.AddSource(ctx, &turbinev2.AddSourceRequest{
				Name: "my-source",
				Plugin: &turbinev2.Plugin{
					Name: "builtin:postgres@1.0.0",
				},
			})
			require.NoError(t, err)

			if test.errMsg == "" {
				test.req.StreamRecords.StreamName = asr.StreamName
			}

			dst, err := s.AddDestination(ctx, &turbinev2.AddDestinationRequest{
				Name: "my-destination",
				Plugin: &turbinev2.Plugin{
					Name: "builtin:postgres@1.0.0",
				},
			})
			require.NoError(t, err)

			if test.errMsg == "" {
				test.req.DestinationID = dst.Id
			}
			_, err = s.WriteRecords(ctx, test.req)
			if test.errMsg != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, test.errMsg)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, s.spec.Streams)
				require.NotEmpty(t, s.spec.Connectors)
				require.Equal(t, s.spec.Streams[0].FromUUID, asr.StreamName)
				require.Equal(t, s.spec.Streams[0].ToUUID, dst.Id)
			}
		})
	}
}

func TestProcessRecords(t *testing.T) {
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

	asr, err := s.AddSource(ctx, &turbinev2.AddSourceRequest{
		Name: "my-source",
		Plugin: &turbinev2.Plugin{
			Name: "builtin:postgres@1.0.0",
		},
	})
	require.NoError(t, err)

	res, err := s.ProcessRecords(ctx, &turbinev2.ProcessRecordsRequest{
		Process: &turbinev2.ProcessRecordsRequest_Process{
			Name: "synchronize",
		},
		StreamRecords: &turbinev2.StreamRecords{
			Records:    []*opencdcv1.Record{},
			StreamName: asr.StreamName,
		},
	})
	require.NoError(t, err)

	require.NotEmpty(t, res)
	require.NotEmpty(t, s.spec.Functions)
	require.Equal(t, s.spec.Streams[0].FromUUID, asr.StreamName)
	require.Equal(t, s.spec.Functions[0].Name, want.Functions[0].Name)
	require.Equal(t, s.spec.Streams[0].ToUUID, res.StreamRecords.StreamName)
}

func TestGetSpec(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		description     string
		populateService func(*SpecBuilderService) *SpecBuilderService
		request         *turbinev2.GetSpecRequest
		want            *ir.DeploymentSpec
		wantErr         error
	}{
		{
			description: "get spec with no function",
			populateService: func(s *SpecBuilderService) *SpecBuilderService {
				s.spec = exampleDeploymentSpec()
				s.spec.Streams = append(s.spec.Streams, ir.StreamSpec{
					UUID:     "1_3",
					FromUUID: "1",
					ToUUID:   "3",
					Name:     "1_3",
				})

				return s
			},
			request: &turbinev2.GetSpecRequest{
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
			populateService: func(s *SpecBuilderService) *SpecBuilderService {
				s.spec = exampleDeploymentSpec()
				s.spec.Streams = append(s.spec.Streams, ir.StreamSpec{
					UUID:     "1_3",
					FromUUID: "1",
					ToUUID:   "3",
					Name:     "1_3",
				})
				return s
			},
			request: &turbinev2.GetSpecRequest{
				Image: "some/image",
			},
			wantErr: fmt.Errorf("cannot set image without defined functions"),
		},
		{
			description: "get spec with function",
			populateService: func(s *SpecBuilderService) *SpecBuilderService {
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
			request: &turbinev2.GetSpecRequest{
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
			populateService: func(s *SpecBuilderService) *SpecBuilderService {
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
			request: &turbinev2.GetSpecRequest{
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
		Connectors: []ir.ConnectorSpec{
			{
				UUID:       "1",
				PluginName: "mongo",
				PluginType: ir.PluginSource,
				PluginConfig: map[string]string{
					"collection": "accounts",
				},
			},
			{
				UUID:       "3",
				PluginName: "postgres",
				PluginType: ir.PluginDestination,
				PluginConfig: map[string]string{
					"collection": "accounts_copy",
					"config":     "value",
				},
			},
		},
		Definition: ir.DefinitionSpec{
			GitSha: "gitsh",
			Metadata: ir.MetadataSpec{
				SpecVersion: "v3",
				Turbine: ir.TurbineSpec{
					Language: ir.GoLang,
					Version:  "10",
				},
			},
		},
	}
}
