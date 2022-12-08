package ir_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/google/go-cmp/cmp"
	"github.com/meroxa/turbine-core/pkg/ir"
)

func Test_DeploymentSpec(t *testing.T) {
	jsonSpec, err := os.ReadFile(path.Join("spectest", "spec.json"))
	if err != nil {
		t.Fatal(err)
	}

	expectedSpec := ir.DeploymentSpec{
		Secrets: map[string]string{
			"key": "valuesecret",
		},
		Connectors: []ir.ConnectorSpec{
			{
				UUID:       "252bc5e1-666e-4985-a12a-42af81a5d2ab",
				Type:       ir.ConnectorSource,
				Resource:   "mypg",
				Collection: "user_activity",
				Config: map[string]interface{}{
					"logical_replication": true,
				},
			},
			{
				UUID:       "dde3bf4e-0848-4579-b05d-7e6dcfae61ea",
				Type:       ir.ConnectorDestination,
				Resource:   "mypg",
				Collection: "user_activity_enriched",
			},
		},
		Functions: []ir.FunctionSpec{
			{
				UUID:  "2ff03fff-6f3e-4f7d-aef8-59c9670bb75d",
				Name:  "user_activity_enriched",
				Image: "ftorres/enrich:9",
			},
		},
		Definition: ir.DefinitionSpec{
			GitSha: "3630e05a-98b7-43a0-aeb0-c9b5b0d4261c",
			Metadata: ir.MetadataSpec{
				Turbine: ir.TurbineSpec{
					Language: ir.GoLang,
					Version:  "0.1.0",
				},
				SpecVersion: "0.2.0",
			},
		},
		Streams: []ir.StreamSpec{
			{
				UUID:     "12345",
				Name:     "my_stream1",
				FromUUID: "252bc5e1-666e-4985-a12a-42af81a5d2ab",
				ToUUID:   "dde3bf4e-0848-4579-b05d-7e6dcfae61ea",
			},
			{
				UUID:     "123456",
				Name:     "my_stream2",
				FromUUID: "dde3bf4e-0848-4579-b05d-7e6dcfae61ea",
				ToUUID:   "2ff03fff-6f3e-4f7d-aef8-59c9670bb75d",
			},
		},
	}

	var deploySpec ir.DeploymentSpec
	if err := json.Unmarshal(jsonSpec, &deploySpec); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(
		deploySpec,
		expectedSpec,
	); diff != "" {
		t.Fatalf("mismatched spec: %s", diff)
	}
}

func Test_ValidateVersion(t *testing.T) {
	testCases := []struct {
		name        string
		specVersion string
		wantError   error
	}{
		{
			name:        "using valid spec version",
			specVersion: ir.LatestSpecVersion,
			wantError:   nil,
		},
		{
			name:        "using invalid spec version",
			specVersion: "0.0.0",
			wantError:   fmt.Errorf("spec version \"0.0.0\" is not a supported. use version %q instead", ir.LatestSpecVersion),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotError := ir.ValidateSpecVersion(tc.specVersion)
			if tc.wantError != nil {
				assert.Equal(t, gotError.Error(), tc.wantError.Error())
			} else {
				assert.NoError(t, gotError)
			}
		})
	}
}

func Test_SetImageForFunctions(t *testing.T) {
	image := "some/image"
	spec := &ir.DeploymentSpec{
		Functions: []ir.FunctionSpec{
			{
				Name: "addition",
			},
			{
				Name: "subtraction",
			},
		},
	}
	spec.SetImageForFunctions(image)

	for _, f := range spec.Functions {
		require.Equal(t, f.Image, image)
	}
}

func Test_ValidateStream(t *testing.T) {
	testCases := []struct {
		spec      ir.DeploymentSpec
		name      string
		wantError error
	}{
		{
			name: "Proper stream ids for from_uuid and to_uuid",
			spec: ir.DeploymentSpec{
				Secrets: map[string]string{
					"a secret": "with value",
				},
				Functions: []ir.FunctionSpec{
					{
						UUID: "252bc5e1-666e-4985-a12a-42af81a5d2ab",
						Name: "addition",
					},
				},
				Connectors: []ir.ConnectorSpec{
					{
						UUID:       "dde3bf4e-0848-4579-b05d-7e6dcfae61ea",
						Collection: "accounts",
						Resource:   "mongo",
						Type:       ir.ConnectorSource,
					},
					{
						UUID:       "2ff03fff-6f3e-4f7d-aef8-59c9670bb75d",
						Collection: "accounts_copy",
						Resource:   "pg",
						Type:       ir.ConnectorDestination,
						Config: map[string]interface{}{
							"config": "value",
						},
					},
				},
				Streams: []ir.StreamSpec{
					{
						UUID:     "12345",
						Name:     "my_stream1",
						FromUUID: "252bc5e1-666e-4985-a12a-42af81a5d2ab",
						ToUUID:   "dde3bf4e-0848-4579-b05d-7e6dcfae61ea",
					},
					{
						UUID:     "123456",
						Name:     "my_stream2",
						FromUUID: "dde3bf4e-0848-4579-b05d-7e6dcfae61ea",
						ToUUID:   "2ff03fff-6f3e-4f7d-aef8-59c9670bb75d",
					},
				},
				Definition: ir.DefinitionSpec{
					GitSha: "gitsh",
					Metadata: ir.MetadataSpec{
						SpecVersion: "0.2.1",
						Turbine: ir.TurbineSpec{
							Language: ir.GoLang,
							Version:  "10",
						},
					},
				},
			},
			wantError: nil,
		},
		{
			name: "Invalid circular stream ids for from_uuid and to_uuid",
			spec: ir.DeploymentSpec{
				Secrets: map[string]string{
					"a secret": "with value",
				},
				Functions: []ir.FunctionSpec{
					{
						UUID: "252bc5e1-666e-4985-a12a-42af81a5d2ab",
						Name: "addition",
					},
				},
				Connectors: []ir.ConnectorSpec{
					{
						UUID:       "dde3bf4e-0848-4579-b05d-7e6dcfae61ea",
						Collection: "accounts",
						Resource:   "mongo",
						Type:       ir.ConnectorSource,
					},
					{
						UUID:       "2ff03fff-6f3e-4f7d-aef8-59c9670bb75d",
						Collection: "accounts_copy",
						Resource:   "pg",
						Type:       ir.ConnectorDestination,
						Config: map[string]interface{}{
							"config": "value",
						},
					},
				},
				Streams: []ir.StreamSpec{
					{
						UUID:     "12345",
						Name:     "my_stream",
						FromUUID: "252bc5e1-666e-4985-a12a-42af81a5d2ab",
						ToUUID:   "252bc5e1-666e-4985-a12a-42af81a5d2ab",
					},
					{
						UUID:     "12345",
						Name:     "my_stream",
						FromUUID: "252bc5e1-666e-4985-a12a-42af81a5d2ab",
						ToUUID:   "dde3bf4e-0848-4579-b05d-7e6dcfae61ea",
					},
				},
				Definition: ir.DefinitionSpec{
					GitSha: "gitsh",
					Metadata: ir.MetadataSpec{
						SpecVersion: "0.2.1",
						Turbine: ir.TurbineSpec{
							Language: ir.GoLang,
							Version:  "10",
						},
					},
				},
			},
			wantError: fmt.Errorf("for stream \"my_stream\" , ids for source (\"252bc5e1-666e-4985-a12a-42af81a5d2ab\") and destination (\"252bc5e1-666e-4985-a12a-42af81a5d2ab\") must be different."),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotError := tc.spec.ValidateStream()
			if tc.wantError != nil {
				assert.Equal(t, gotError.Error(), tc.wantError.Error())
			} else {
				assert.NoError(t, gotError)
			}
		})
	}
}

func Test_MarshalUnmarshal(t *testing.T) {
	spec := &ir.DeploymentSpec{
		Secrets: map[string]string{
			"a secret": "with value",
		},
		Functions: []ir.FunctionSpec{
			{
				Name: "addition",
			},
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
		Streams: []ir.StreamSpec{
			{
				UUID:     "12345",
				Name:     "my_stream",
				FromUUID: "252bc5e1-666e-4985-a12a-42af81a5d2ab",
				ToUUID:   "dde3bf4e-0848-4579-b05d-7e6dcfae61ea",
			},
		},
		Definition: ir.DefinitionSpec{
			GitSha: "gitsh",
			Metadata: ir.MetadataSpec{
				SpecVersion: "0.2.1",
				Turbine: ir.TurbineSpec{
					Language: ir.GoLang,
					Version:  "10",
				},
			},
		},
	}
	b, err := spec.Marshal()
	require.NoError(t, err)

	got, err := ir.Unmarshal(b)
	require.NoError(t, err)

	require.Equal(t, spec, got)
}
