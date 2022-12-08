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
				ID:         "1",
				Type:       ir.ConnectorSource,
				Resource:   "mypg",
				Collection: "user_activity",
				Config: map[string]interface{}{
					"logical_replication": true,
				},
			},
			{
				ID:         "2",
				Type:       ir.ConnectorDestination,
				Resource:   "mypg",
				Collection: "user_activity_enriched",
			},
		},
		Functions: []ir.FunctionSpec{
			{
				ID:    "3",
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
				ID:     "12345",
				Name:   "my_stream1",
				FromID: "1",
				ToID:   "2",
			},
			{
				ID:     "123456",
				Name:   "my_stream2",
				FromID: "2",
				ToID:   "3",
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

func Test_ValidateStreamIDs(t *testing.T) {
	testCases := []struct {
		spec      ir.DeploymentSpec
		name      string
		wantError error
	}{
		{
			name: "Proper stream ids for from_id and to_id",
			spec: ir.DeploymentSpec{
				Secrets: map[string]string{
					"a secret": "with value",
				},
				Functions: []ir.FunctionSpec{
					{
						ID:   "1",
						Name: "addition",
					},
				},
				Connectors: []ir.ConnectorSpec{
					{
						ID:         "2",
						Collection: "accounts",
						Resource:   "mongo",
						Type:       ir.ConnectorSource,
					},
					{
						ID:         "3",
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
						ID:     "12345",
						Name:   "my_stream1",
						FromID: "1",
						ToID:   "2",
					},
					{
						ID:     "123456",
						Name:   "my_stream2",
						FromID: "2",
						ToID:   "3",
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
			name: "Invalid circular stream ids for from_id and to_id",
			spec: ir.DeploymentSpec{
				Secrets: map[string]string{
					"a secret": "with value",
				},
				Functions: []ir.FunctionSpec{
					{
						ID:   "1",
						Name: "addition",
					},
				},
				Connectors: []ir.ConnectorSpec{
					{
						ID:         "2",
						Collection: "accounts",
						Resource:   "mongo",
						Type:       ir.ConnectorSource,
					},
					{
						ID:         "3",
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
						ID:     "12345",
						Name:   "my_stream",
						FromID: "1",
						ToID:   "1",
					},
					{
						ID:     "12345",
						Name:   "my_stream",
						FromID: "1",
						ToID:   "2",
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
			wantError: fmt.Errorf("for stream \"my_stream\" , ids for source (\"1\") and destination (\"1\") must be different."),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotError := tc.spec.ValidateStreamIDs()
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
				ID:     "12345",
				Name:   "my_stream",
				FromID: "1",
				ToID:   "2",
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
