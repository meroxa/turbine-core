package ir_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/heimdalr/dag"
	"github.com/meroxa/turbine-core/v2/pkg/ir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeploymentSpec_BuildDAG_UnsupportedSpec(t *testing.T) {
	jsonSpec, err := os.ReadFile(path.Join("v3", "spectest", "spec_unsupported.json"))
	if err != nil {
		t.Fatal(err)
	}

	var spec ir.DeploymentSpec
	if err := json.Unmarshal(jsonSpec, &spec); err != nil {
		t.Fatal(err)
	}

	_, err = spec.BuildDAG()
	assert.ErrorContains(t, err, "spec version \"0.0.0\" is invalid, supported versions: v3")
}

func TestDeploymentSpec_BuildDAG_EmptySpec(t *testing.T) {
	jsonSpec, err := os.ReadFile(path.Join("v3", "spectest", "spec_empty_ver.json"))
	if err != nil {
		t.Fatal(err)
	}

	var spec ir.DeploymentSpec
	if err := json.Unmarshal(jsonSpec, &spec); err != nil {
		t.Fatal(err)
	}

	_, err = spec.BuildDAG()
	assert.ErrorContains(t, err, "spec version \"\" is invalid, supported versions: v3")
}

func Test_DeploymentSpec(t *testing.T) {
	jsonSpec, err := os.ReadFile(path.Join("v3", "spectest", "spec.json"))
	if err != nil {
		t.Fatal(err)
	}

	expectedSpec := &ir.DeploymentSpec{
		Connectors: []ir.ConnectorSpec{
			{
				UUID:       "252bc5e1-666e-4985-a12a-42af81a5d2ab",
				PluginType: ir.PluginSource,
				PluginName: "postgres",
				PluginConfig: map[string]string{
					"collection":          "user_activity",
					"logical_replication": "true",
				},
			},
			{
				UUID:       "dde3bf4e-0848-4579-b05d-7e6dcfae61ea",
				PluginType: ir.PluginDestination,
				PluginName: "postgres",
				PluginConfig: map[string]string{
					"collection": "user_activity_enriched",
				},
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
				SpecVersion: "v3",
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

	deploySpec := &ir.DeploymentSpec{}
	if err := json.Unmarshal(jsonSpec, deploySpec); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, deploySpec, expectedSpec)
}

func Test_ValidateVersion(t *testing.T) {
	testCases := []struct {
		name         string
		specVersions []string
		wantError    error
	}{
		{
			name:         "using valid spec version",
			specVersions: []string{"v3"},
			wantError:    nil,
		},
		{
			name:         "using invalid spec version",
			specVersions: []string{"0.0.0"},
			wantError:    fmt.Errorf("spec version \"0.0.0\" is invalid, supported versions: v3"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, v := range tc.specVersions {
				err := ir.ValidateSpecVersion(v)
				if tc.wantError != nil {
					assert.Equal(t, err.Error(), tc.wantError.Error())
				} else {
					assert.NoError(t, err)
				}
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

func Test_MarshalUnmarshal(t *testing.T) {
	spec := &ir.DeploymentSpec{
		Functions: []ir.FunctionSpec{
			{
				UUID: "3",
				Name: "addition",
			},
		},
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
				UUID:       "2",
				PluginName: "pg",
				PluginType: ir.PluginDestination,
				PluginConfig: map[string]string{
					"collection": "accounts_copy",
					"config":     "value",
				},
			},
		},
		Streams: []ir.StreamSpec{
			{
				UUID:     "12345",
				Name:     "my_stream",
				FromUUID: "1",
				ToUUID:   "2",
			},
			{
				UUID:     "12345",
				Name:     "my_stream2",
				FromUUID: "2",
				ToUUID:   "3",
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
	b, err := spec.Marshal()
	require.NoError(t, err)

	got, err := ir.Unmarshal(b)
	require.NoError(t, err)

	require.Equal(t, spec, got)
}

func Test_AllowMultipleSources(t *testing.T) {
	var spec ir.DeploymentSpec
	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			PluginName: "mongo",
			PluginType: ir.PluginSource,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	)
	require.NoError(t, err)
	err = spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "2",
			PluginName: "mongo",
			PluginType: ir.PluginSource,
			PluginConfig: map[string]string{
				"collection": "accounts2",
				"config":     "value",
			},
		},
	)
	require.NoError(t, err)
}

func Test_EnsureNonDuplicateSources(t *testing.T) {
	var spec ir.DeploymentSpec
	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			PluginName: "mongo",
			PluginType: ir.PluginSource,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	)
	require.NoError(t, err)

	duplicate := &ir.ConnectorSpec{
		UUID:       "1",
		PluginName: "mongo",
		PluginType: ir.PluginSource,
		PluginConfig: map[string]string{
			"collection": "accounts2",
			"config":     "value",
		},
	}

	err = spec.AddSource(duplicate)
	require.EqualError(t, err, fmt.Errorf("the id '%s' is already known", duplicate.UUID).Error())
}

func Test_BadStream(t *testing.T) {
	var spec ir.DeploymentSpec
	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			PluginName: "mongo",
			PluginType: ir.PluginSource,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	)
	require.NoError(t, err)
	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "1_2",
			Name:     "my_stream1",
			FromUUID: "1",
			ToUUID:   "2",
		},
	)
	require.Error(t, err)
	require.Equal(t, err.Error(), "destination 2 does not exist")
}

func Test_WrongSourceConnector(t *testing.T) {
	var spec ir.DeploymentSpec
	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			PluginName: "mongo",
			PluginType: ir.PluginDestination,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	)
	require.Error(t, err)
	require.Equal(t, err.Error(), "not a source connector")
}

func Test_WrongDestinationConnector(t *testing.T) {
	var spec ir.DeploymentSpec
	err := spec.AddDestination(
		&ir.ConnectorSpec{
			UUID:       "1",
			PluginName: "mongo",
			PluginType: ir.PluginSource,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	)
	require.Error(t, err)
	require.Equal(t, err.Error(), "not a destination connector")
}

// Scenario 1 - Simple DAG
// source → fn -> dest
// ( src_con ) → (stream) → (function) → (stream) → (dest1) .
func Test_Scenario1(t *testing.T) {
	var spec ir.DeploymentSpec
	spec.Definition.Metadata.SpecVersion = ir.SpecVersion_v3

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			PluginName: "mongo",
			PluginType: ir.PluginSource,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddFunction(
		&ir.FunctionSpec{
			UUID: "2",
			Name: "addition",
		},
	)
	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "1_2",
			Name:     "my_stream1",
			FromUUID: "1",
			ToUUID:   "2",
		},
	)

	require.NoError(t, err)

	err = spec.AddDestination(
		&ir.ConnectorSpec{
			UUID:       "3",
			PluginName: "pg",
			PluginType: ir.PluginDestination,
			PluginConfig: map[string]string{
				"collection": "accounts_copy",
				"config":     "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "2_3",
			Name:     "my_stream2",
			FromUUID: "2",
			ToUUID:   "3",
		},
	)
	require.NoError(t, err)

	_, err = spec.BuildDAG()
	require.NoError(t, err)
}

// Scenario 2 - DAG with two Destinations from 1 function
// source → fn -> dest[n]
// ( src_con ) → (stream) → (function) → (stream) → (dest1)
//
//			↓
//	    (stream) → (dest2)
func Test_DAGScenario2(t *testing.T) {
	var spec ir.DeploymentSpec
	spec.Definition.Metadata.SpecVersion = ir.SpecVersion_v3

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			PluginName: "mongo",
			PluginType: ir.PluginSource,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddFunction(
		&ir.FunctionSpec{
			UUID: "2",
			Name: "addition",
		},
	)
	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "1_2",
			Name:     "my_stream1",
			FromUUID: "1",
			ToUUID:   "2",
		},
	)

	require.NoError(t, err)

	err = spec.AddDestination(
		&ir.ConnectorSpec{
			UUID:       "3",
			PluginName: "pg",
			PluginType: ir.PluginDestination,
			PluginConfig: map[string]string{
				"collection": "accounts_copy",
				"config":     "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "2_3",
			Name:     "my_stream2",
			FromUUID: "2",
			ToUUID:   "3",
		},
	)

	require.NoError(t, err)

	err = spec.AddDestination(
		&ir.ConnectorSpec{
			UUID:       "4",
			PluginName: "pg",
			PluginType: ir.PluginDestination,
			PluginConfig: map[string]string{
				"collection": "accounts_copy",
				"config":     "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "2_4",
			Name:     "my_stream4",
			FromUUID: "2",
			ToUUID:   "4",
		},
	)
	require.NoError(t, err)

	_, err = spec.BuildDAG()
	require.NoError(t, err)
}

// Scenario 3 - Not a DAG, trying to write from one function back to the other creates a loop
// source → (fn) → dest[n]
//
//	                                 ↓   ←      ←      ←      ↑
//	   ( src_con ) → (stream) → (function) → (stream) →  (function)
//										↓
//							    	(stream) →	(dest2)
func Test_DAGScenario3(t *testing.T) {
	var spec ir.DeploymentSpec

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			PluginName: "mongo",
			PluginType: ir.PluginSource,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddFunction(
		&ir.FunctionSpec{
			UUID: "2",
			Name: "addition",
		},
	)
	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "1_2",
			Name:     "my_stream1",
			FromUUID: "1",
			ToUUID:   "2",
		},
	)

	require.NoError(t, err)

	err = spec.AddDestination(
		&ir.ConnectorSpec{
			UUID:       "3",
			PluginName: "pg",
			PluginType: ir.PluginDestination,
			PluginConfig: map[string]string{
				"collection": "accounts_copy",
				"config":     "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "2_3",
			Name:     "my_stream2",
			FromUUID: "2",
			ToUUID:   "3",
		},
	)

	require.NoError(t, err)

	err = spec.AddFunction(
		&ir.FunctionSpec{
			UUID: "4",
			Name: "addition",
		},
	)
	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "2_4",
			Name:     "my_stream3",
			FromUUID: "2",
			ToUUID:   "4",
		},
	)

	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "4_2",
			Name:     "my_stream3",
			FromUUID: "4",
			ToUUID:   "2",
		},
	)

	require.Error(t, err)
	assert.Equal(t, err.Error(), "edge between '4' and '2' would create a loop")
}

// Scenario 4 - Not acyclic, trying to write from one function back to source
//
//	   ( src_con ) → (stream) → (function 1)
//				↑					    ↓
//			    ←  ←  ← (stream) ← ← ← ←
func Test_DAGScenario4(t *testing.T) {
	var spec ir.DeploymentSpec

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			PluginName: "mongo",
			PluginType: ir.PluginSource,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddFunction(
		&ir.FunctionSpec{
			UUID: "2",
			Name: "addition",
		},
	)
	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "1_2",
			Name:     "my_stream1",
			FromUUID: "1",
			ToUUID:   "2",
		},
	)

	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "2_1",
			Name:     "my_stream2",
			FromUUID: "2",
			ToUUID:   "1",
		},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "would create a loop")
}

// Scenario 5 - DAG, multuple functions, 1 destination
//
//	source → (fn) → (fn)… → dest
//	  ( src_con ) → (stream) → (function 1) → (stream) → (function2)  → (dest1)
func Test_DAGScenario5(t *testing.T) {
	spec := ir.DeploymentSpec{
		Definition: ir.DefinitionSpec{
			Metadata: ir.MetadataSpec{SpecVersion: ir.LatestSpecVersion},
		},
	}

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			PluginName: "mongo",
			PluginType: ir.PluginSource,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddFunction(
		&ir.FunctionSpec{
			UUID: "2",
			Name: "addition",
		},
	)
	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "1_2",
			Name:     "my_stream1",
			FromUUID: "1",
			ToUUID:   "2",
		},
	)

	require.NoError(t, err)

	err = spec.AddFunction(
		&ir.FunctionSpec{
			UUID: "3",
			Name: "addition_more2",
		},
	)
	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "2_3",
			Name:     "my_stream4",
			FromUUID: "2",
			ToUUID:   "3",
		},
	)
	require.NoError(t, err)

	err = spec.AddDestination(
		&ir.ConnectorSpec{
			UUID:       "4",
			PluginName: "pg",
			PluginType: ir.PluginDestination,
			PluginConfig: map[string]string{
				"collection": "accounts_copy",
				"config":     "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "3_4",
			Name:     "my_stream4",
			FromUUID: "3",
			ToUUID:   "4",
		},
	)
	require.NoError(t, err)

	_, err = spec.BuildDAG()
	require.NoError(t, err)
}

// Scenario 6 - DAG with many functions and destinations
// source → (fn) → (fn)… → dest[n]
//
//	   ( src_con ) → (stream) → (function 1) → (stream) → (function2)  → (dest1)
//									    ↓
//									 (stream)  → (dest2)
//									    ↓
//									(function 3)
//										↓
//									(stream)
//	  								↓
//									(dest 3)
func Test_DAGScenario6(t *testing.T) {
	spec := ir.DeploymentSpec{
		Definition: ir.DefinitionSpec{
			Metadata: ir.MetadataSpec{SpecVersion: ir.LatestSpecVersion},
		},
	}

	require.NoError(t, spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			PluginName: "mongo",
			PluginType: ir.PluginSource,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	))

	require.NoError(t, spec.AddFunction(
		&ir.FunctionSpec{
			UUID: "2",
			Name: "addition_first_function",
		},
	))

	require.NoError(t, spec.AddStream(
		&ir.StreamSpec{
			UUID:     "1_2",
			Name:     "my_stream1",
			FromUUID: "1",
			ToUUID:   "2",
		},
	))

	require.NoError(t, spec.AddFunction(
		&ir.FunctionSpec{
			UUID: "3",
			Name: "subtraction_second_function",
		},
	))

	require.NoError(t, spec.AddStream(
		&ir.StreamSpec{
			UUID:     "2_3",
			Name:     "my_stream4",
			FromUUID: "2",
			ToUUID:   "3",
		},
	))

	require.NoError(t, spec.AddDestination(
		&ir.ConnectorSpec{
			UUID:       "4",
			PluginName: "pg",
			PluginType: ir.PluginDestination,
			PluginConfig: map[string]string{
				"collection": "accounts_copy",
				"config":     "value",
			},
		},
	))

	require.NoError(t, spec.AddStream(
		&ir.StreamSpec{
			UUID:     "3_4",
			Name:     "my_stream4",
			FromUUID: "3",
			ToUUID:   "4",
		},
	))

	require.NoError(t, spec.AddDestination(
		&ir.ConnectorSpec{
			UUID:       "5",
			PluginName: "pg",
			PluginType: ir.PluginDestination,
			PluginConfig: map[string]string{
				"collection": "accounts_copy_2",
				"config":     "value",
			},
		},
	))

	require.NoError(t, spec.AddStream(
		&ir.StreamSpec{
			UUID:     "2_5",
			Name:     "my_stream5",
			FromUUID: "2",
			ToUUID:   "5",
		},
	))

	require.NoError(t, spec.AddFunction(
		&ir.FunctionSpec{
			UUID: "6",
			Name: "multiplication_third_function",
		},
	))

	require.NoError(t, spec.AddStream(
		&ir.StreamSpec{
			UUID:     "2_6",
			Name:     "my_stream5",
			FromUUID: "2",
			ToUUID:   "6",
		},
	))

	require.NoError(t, spec.AddDestination(
		&ir.ConnectorSpec{
			UUID:       "7",
			PluginName: "pg",
			PluginType: ir.PluginDestination,
			PluginConfig: map[string]string{
				"collection": "accounts_copy_3",
				"config":     "value",
			},
		},
	))

	require.NoError(t, spec.AddStream(
		&ir.StreamSpec{
			UUID:     "6_7",
			Name:     "my_stream6",
			FromUUID: "6",
			ToUUID:   "7",
		},
	))
	_, err := spec.BuildDAG()
	require.NoError(t, err)
}

// Scenario 7 - DAG with destination from 1 function and with source data going to second destination
// source → dest[0] | (fn)→ dest[1]
// ( src_con ) → (stream) → (function) → (stream) → (dest1)
//
//		↓
//	(stream) → (dest2)
func Test_DAGScenario7(t *testing.T) {
	var spec ir.DeploymentSpec
	spec.Definition.Metadata.SpecVersion = ir.SpecVersion_v3

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			PluginName: "mongo",
			PluginType: ir.PluginSource,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddFunction(
		&ir.FunctionSpec{
			UUID: "2",
			Name: "addition",
		},
	)
	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "1_2",
			Name:     "my_stream1",
			FromUUID: "1",
			ToUUID:   "2",
		},
	)

	require.NoError(t, err)

	err = spec.AddDestination(
		&ir.ConnectorSpec{
			UUID:       "3",
			PluginName: "pg",
			PluginType: ir.PluginDestination,
			PluginConfig: map[string]string{
				"collection": "accounts_copy",
				"config":     "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "2_3",
			Name:     "my_stream2",
			FromUUID: "2",
			ToUUID:   "3",
		},
	)

	require.NoError(t, err)

	err = spec.AddDestination(
		&ir.ConnectorSpec{
			UUID:       "4",
			PluginName: "pg",
			PluginType: ir.PluginDestination,
			PluginConfig: map[string]string{
				"collection": "accounts_copy",
				"config":     "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "1_4",
			Name:     "my_stream4",
			FromUUID: "1",
			ToUUID:   "4",
		},
	)
	require.NoError(t, err)

	_, err = spec.BuildDAG()
	require.NoError(t, err)
}

// Scenario 8 - Just source to function, no destination
// source → fn
// ( src_con ) → (stream) → (func)
func Test_Scenario8(t *testing.T) {
	var spec ir.DeploymentSpec
	spec.Definition.Metadata.SpecVersion = ir.SpecVersion_v3

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			PluginName: "mongo",
			PluginType: ir.PluginSource,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddFunction(
		&ir.FunctionSpec{
			UUID: "2",
			Name: "addition",
		},
	)
	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "1_2",
			Name:     "my_stream1",
			FromUUID: "1",
			ToUUID:   "2",
		},
	)
	require.NoError(t, err)
	_, err = spec.BuildDAG()
	require.NoError(t, err)
}

// Scenario 9 - Just source to destination, no function
// source → dest
// ( src_con ) → (stream) → (destination)
func Test_Scenario9(t *testing.T) {
	var spec ir.DeploymentSpec
	spec.Definition.Metadata.SpecVersion = ir.SpecVersion_v3

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			PluginName: "mongo",
			PluginType: ir.PluginSource,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddDestination(
		&ir.ConnectorSpec{
			UUID:       "2",
			PluginName: "mongo",
			PluginType: ir.PluginDestination,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	)
	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "1_2",
			Name:     "my_stream1",
			FromUUID: "1",
			ToUUID:   "2",
		},
	)
	require.NoError(t, err)

	_, err = spec.BuildDAG()
	require.NoError(t, err)
}

// Scenario 10 - Disconnected Graph
// src -> fn[0]
// fn[1] -> dst
func Test_Scenario10(t *testing.T) {
	spec := ir.DeploymentSpec{
		Definition: ir.DefinitionSpec{
			Metadata: ir.MetadataSpec{SpecVersion: ir.LatestSpecVersion},
		},
	}

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			PluginName: "mongo",
			PluginType: ir.PluginSource,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddFunction(
		&ir.FunctionSpec{
			UUID:  "2",
			Name:  "function",
			Image: "test",
		},
	)

	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "1_2",
			Name:     "my_stream1",
			FromUUID: "1",
			ToUUID:   "2",
		},
	)

	require.NoError(t, err)

	err = spec.AddFunction(
		&ir.FunctionSpec{
			UUID:  "3",
			Name:  "function2",
			Image: "test",
		},
	)

	require.NoError(t, err)

	err = spec.AddDestination(
		&ir.ConnectorSpec{
			UUID:       "4",
			PluginName: "mongo",
			PluginType: ir.PluginDestination,
			PluginConfig: map[string]string{
				"collection": "accounts",
				"config":     "value",
			},
		},
	)
	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "3_4",
			Name:     "my_stream1",
			FromUUID: "3",
			ToUUID:   "4",
		},
	)
	require.NoError(t, err)
	dag, err := spec.BuildDAG()
	require.NoError(t, err)

	err = spec.ValidateDAG(dag)
	require.Error(t, err)
	assert.Equal(t, err.Error(), "invalid DAG, too many sources")
}

func Test_ValidateDAG(t *testing.T) {
	testCases := []struct {
		name      string
		setup     func(t *testing.T) *dag.DAG
		wantError error
	}{
		{
			name: "empty DAG",
			setup: func(t *testing.T) *dag.DAG {
				t.Helper()

				spec := ir.DeploymentSpec{
					Definition: ir.DefinitionSpec{
						Metadata: ir.MetadataSpec{
							SpecVersion: ir.SpecVersion_v3,
						},
					},
				}

				dag, err := spec.BuildDAG()
				require.NoError(t, err)

				return dag
			},
			wantError: fmt.Errorf("invalid DAG, no sources found"),
		},
		{
			name: "too many sources",
			setup: func(t *testing.T) *dag.DAG {
				t.Helper()

				spec := ir.DeploymentSpec{
					Definition: ir.DefinitionSpec{
						Metadata: ir.MetadataSpec{
							SpecVersion: ir.SpecVersion_v3,
						},
					},
					Connectors: []ir.ConnectorSpec{
						{
							UUID:       uuid.New().String(),
							PluginType: ir.PluginSource,
						},
						{
							UUID:       uuid.New().String(),
							PluginType: ir.PluginSource,
						},
					},
				}

				dag, err := spec.BuildDAG()
				require.NoError(t, err)

				return dag
			},
			wantError: fmt.Errorf("invalid DAG, too many sources"),
		},
		{
			name: "only one source",
			setup: func(t *testing.T) *dag.DAG {
				t.Helper()

				spec := ir.DeploymentSpec{
					Definition: ir.DefinitionSpec{
						Metadata: ir.MetadataSpec{
							SpecVersion: ir.SpecVersion_v3,
						},
					},
					Connectors: []ir.ConnectorSpec{
						{
							PluginType: ir.PluginSource,
						},
					},
				}

				dag, err := spec.BuildDAG()
				require.NoError(t, err)

				return dag
			},
			wantError: fmt.Errorf("invalid DAG, there has to be at least one source, at most one function, and zero or more destinations"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spec := &ir.DeploymentSpec{}

			err := spec.ValidateDAG(tc.setup(t))
			if tc.wantError != nil {
				require.Error(t, err)
				require.Equal(t, tc.wantError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
