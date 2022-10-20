package ir_test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"

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
				Type:       ir.ConnectorSource,
				Resource:   "mypg",
				Collection: "user_activity",
				Config: map[string]interface{}{
					"logical_replication": true,
				},
			},
			{
				Type:       ir.ConnectorDestination,
				Resource:   "mypg",
				Collection: "user_activity_enriched",
			},
		},
		Functions: []ir.FunctionSpec{
			{
				Name:  "user_activity_enriched",
				Image: "ftorres/enrich:9",
				EnvVars: map[string]interface{}{
					"CLEARBIT_API_KEY": "token-1",
				},
			},
		},
		Definition: ir.DefinitionSpec{
			GitSha: "3630e05a-98b7-43a0-aeb0-c9b5b0d4261c",
			Metadata: ir.MetadataSpec{
				Turbine: ir.TurbineSpec{
					Language: ir.GoLang,
					Version:  "0.1.0",
				},
				SpecVersion: "0.1.1",
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
