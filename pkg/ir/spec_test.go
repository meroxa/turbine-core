package ir_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/meroxa/turbine-core/pkg/ir"
)

func TestDeploymentSpec_BuildDAG_UnsupportedUpgrade(t *testing.T) {
	jsonSpec, err := os.ReadFile(path.Join("spectest", "0.0.0", "spec.json"))
	if err != nil {
		t.Fatal(err)
	}

	var spec ir.DeploymentSpec
	if err := json.Unmarshal(jsonSpec, &spec); err != nil {
		t.Fatal(err)
	}

	_, err = spec.BuildDAG()
	assert.ErrorContains(t, err, fmt.Sprintf("unsupported upgrade from spec version \"0.0.0\" to %q", ir.LatestSpecVersion))
}

func TestDeploymentSpec_BuildDAG_0_1_1(t *testing.T) {
	jsonSpec, err := os.ReadFile(path.Join("spectest", "0.1.1", "spec.json"))
	if err != nil {
		t.Fatal(err)
	}

	var spec ir.DeploymentSpec
	if err := json.Unmarshal(jsonSpec, &spec); err != nil {
		t.Fatal(err)
	}

	dag, err := spec.BuildDAG()
	require.NoError(t, err)

	var fnUUID, destUUID string

	// Check root is a connector source
	roots := dag.GetRoots()
	assert.Equal(t, len(roots), 1)
	for _, s := range roots {
		connector, ok := s.(*ir.ConnectorSpec)
		if !ok {
			t.Fatalf("root edge is not a connector")
		}
		assert.Equal(t, connector.Type, ir.ConnectorSource)
	}

	// Check its only leaf is a connector destination
	leaves := dag.GetLeaves()
	assert.Equal(t, len(leaves), 1)
	for _, s := range leaves {
		connector, ok := s.(*ir.ConnectorSpec)
		if !ok {
			t.Fatalf("leaf edge is not a connector")
		}
		destUUID = connector.UUID
		assert.Equal(t, connector.Type, ir.ConnectorDestination)
	}

	// Check function connects both source and destination

	// From destination connector
	fnEdges, err := dag.GetParents(destUUID)
	assert.NoError(t, err)

	for _, fn := range fnEdges {
		function, ok := fn.(*ir.FunctionSpec)
		if !ok {
			t.Fatalf("edge is not a not a function")
		}
		fnUUID = function.UUID
	}

	// From the function itself checks its parent is a connector source
	srcEdges, err := dag.GetParents(fnUUID)
	assert.NoError(t, err)

	for _, src := range srcEdges {
		connector, ok := src.(*ir.ConnectorSpec)
		if !ok {
			t.Fatalf("edge is not a not a connector")
		}
		assert.Equal(t, connector.Type, ir.ConnectorSource)
	}

	assert.Equal(t, len(dag.GetVertices()), 3)

	// Number of edges created
	assert.Equal(t, dag.GetSize(), 2)
}

func Test_DeploymentSpec(t *testing.T) {
	jsonSpec, err := os.ReadFile(path.Join("spectest", "0.2.0", "spec.json"))
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
			specVersions: []string{"0.1.1", "0.2.0"},
			wantError:    nil,
		},
		{},
		{
			name:         "using invalid spec version",
			specVersions: []string{"0.0.0"},
			wantError:    fmt.Errorf("spec version \"0.0.0\" is invalid, supported versions: 0.1.1, 0.2.0"),
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
		Secrets: map[string]string{
			"a secret": "with value",
		},
		Functions: []ir.FunctionSpec{
			{
				UUID: "3",
				Name: "addition",
			},
		},
		Connectors: []ir.ConnectorSpec{
			{
				UUID:       "1",
				Collection: "accounts",
				Resource:   "mongo",
				Type:       ir.ConnectorSource,
			},
			{
				UUID:       "2",
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
				SpecVersion: "0.2.0",
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

func Test_EnsureSingleSource(t *testing.T) {
	var spec ir.DeploymentSpec
	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			Collection: "accounts",
			Resource:   "mongo",
			Type:       ir.ConnectorSource,
			Config: map[string]interface{}{
				"config": "value",
			},
		},
	)
	require.NoError(t, err)
	err = spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			Collection: "accounts2",
			Resource:   "mongo",
			Type:       ir.ConnectorSource,
			Config: map[string]interface{}{
				"config": "value",
			},
		},
	)
	require.EqualError(t, err, fmt.Errorf("can only add one source connector per application").Error())
}

func Test_BadStream(t *testing.T) {
	var spec ir.DeploymentSpec
	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			Collection: "accounts",
			Resource:   "mongo",
			Type:       ir.ConnectorSource,
			Config: map[string]interface{}{
				"config": "value",
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
			Collection: "accounts",
			Resource:   "mongo",
			Type:       ir.ConnectorDestination,
			Config: map[string]interface{}{
				"config": "value",
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
			Collection: "accounts",
			Resource:   "mongo",
			Type:       ir.ConnectorSource,
			Config: map[string]interface{}{
				"config": "value",
			},
		},
	)
	require.Error(t, err)
	require.Equal(t, err.Error(), "not a destination connector")
}

// Scenario 1 - Simple DAG
// source → fn -> dest
// ( src_con ) → (stream) → (function) → (stream) → (dest1)
func Test_Scenario1(t *testing.T) {
	var spec ir.DeploymentSpec
	spec.Definition.Metadata.SpecVersion = ir.SpecVersion_0_2_0

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			Collection: "accounts",
			Resource:   "mongo",
			Type:       ir.ConnectorSource,
			Config: map[string]interface{}{
				"config": "value",
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
			Collection: "accounts_copy",
			Resource:   "pg",
			Type:       ir.ConnectorDestination,
			Config: map[string]interface{}{
				"config": "value",
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
	spec.Definition.Metadata.SpecVersion = ir.SpecVersion_0_2_0

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			Collection: "accounts",
			Resource:   "mongo",
			Type:       ir.ConnectorSource,
			Config: map[string]interface{}{
				"config": "value",
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
			Collection: "accounts_copy",
			Resource:   "pg",
			Type:       ir.ConnectorDestination,
			Config: map[string]interface{}{
				"config": "value",
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
			Collection: "accounts_copy",
			Resource:   "pg",
			Type:       ir.ConnectorDestination,
			Config: map[string]interface{}{
				"config": "value",
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
			Collection: "accounts",
			Resource:   "mongo",
			Type:       ir.ConnectorSource,
			Config: map[string]interface{}{
				"config": "value",
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
			Collection: "accounts_copy",
			Resource:   "pg",
			Type:       ir.ConnectorDestination,
			Config: map[string]interface{}{
				"config": "value",
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
			Collection: "accounts",
			Resource:   "mongo",
			Type:       ir.ConnectorSource,
			Config: map[string]interface{}{
				"config": "value",
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
	var spec ir.DeploymentSpec

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			Collection: "accounts",
			Resource:   "mongo",
			Type:       ir.ConnectorSource,
			Config: map[string]interface{}{
				"config": "value",
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
			Collection: "accounts_copy",
			Resource:   "pg",
			Type:       ir.ConnectorDestination,
			Config: map[string]interface{}{
				"config": "value",
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
	var spec ir.DeploymentSpec

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			Collection: "accounts",
			Resource:   "mongo",
			Type:       ir.ConnectorSource,
			Config: map[string]interface{}{
				"config": "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddFunction(
		&ir.FunctionSpec{
			UUID: "2",
			Name: "addition_first_function",
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
			Name: "subtraction_second_function",
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
			Collection: "accounts_copy",
			Resource:   "pg",
			Type:       ir.ConnectorDestination,
			Config: map[string]interface{}{
				"config": "value",
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

	err = spec.AddDestination(
		&ir.ConnectorSpec{
			UUID:       "5",
			Collection: "accounts_copy_2",
			Resource:   "pg",
			Type:       ir.ConnectorDestination,
			Config: map[string]interface{}{
				"config": "value",
			},
		},
	)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "2_5",
			Name:     "my_stream5",
			FromUUID: "2",
			ToUUID:   "5",
		},
	)

	require.NoError(t, err)

	err = spec.AddFunction(
		&ir.FunctionSpec{
			UUID: "6",
			Name: "multiplication_third_function",
		},
	)
	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "2_6",
			Name:     "my_stream5",
			FromUUID: "2",
			ToUUID:   "6",
		},
	)
	require.NoError(t, err)

	err = spec.AddDestination(
		&ir.ConnectorSpec{
			UUID:       "7",
			Collection: "accounts_copy_3",
			Resource:   "pg",
			Type:       ir.ConnectorDestination,
			Config: map[string]interface{}{
				"config": "value",
			},
		},
	)
	require.NoError(t, err)

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "6_7",
			Name:     "my_stream6",
			FromUUID: "6",
			ToUUID:   "7",
		},
	)
	require.NoError(t, err)

	_, err = spec.BuildDAG()
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
	spec.Definition.Metadata.SpecVersion = ir.SpecVersion_0_2_0

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			Collection: "accounts",
			Resource:   "mongo",
			Type:       ir.ConnectorSource,
			Config: map[string]interface{}{
				"config": "value",
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
			Collection: "accounts_copy",
			Resource:   "pg",
			Type:       ir.ConnectorDestination,
			Config: map[string]interface{}{
				"config": "value",
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
			Collection: "accounts_copy",
			Resource:   "pg",
			Type:       ir.ConnectorDestination,
			Config: map[string]interface{}{
				"config": "value",
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
	spec.Definition.Metadata.SpecVersion = ir.SpecVersion_0_2_0

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			Collection: "accounts",
			Resource:   "mongo",
			Type:       ir.ConnectorSource,
			Config: map[string]interface{}{
				"config": "value",
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
	spec.Definition.Metadata.SpecVersion = ir.SpecVersion_0_2_0

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			Collection: "accounts",
			Resource:   "mongo",
			Type:       ir.ConnectorSource,
			Config: map[string]interface{}{
				"config": "value",
			},
		},
	)

	require.NoError(t, err)

	err = spec.AddDestination(
		&ir.ConnectorSpec{
			UUID:       "2",
			Collection: "accounts",
			Resource:   "mongo",
			Type:       ir.ConnectorDestination,
			Config: map[string]interface{}{
				"config": "value",
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
	var spec ir.DeploymentSpec

	err := spec.AddSource(
		&ir.ConnectorSpec{
			UUID:       "1",
			Collection: "accounts",
			Resource:   "mongo",
			Type:       ir.ConnectorSource,
			Config: map[string]interface{}{
				"config": "value",
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
			Collection: "accounts",
			Resource:   "mongo",
			Type:       ir.ConnectorDestination,
			Config: map[string]interface{}{
				"config": "value",
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
	_, err = spec.BuildDAG()
	require.Error(t, err)
	assert.Equal(t, err.Error(), "invalid DAG, too many sources")
}

func Test_ValidateDAG(t *testing.T) {
	testCases := []struct {
		name      string
		spec      ir.DeploymentSpec
		wantError error
	}{
		{
			name: "empty DAG",
			spec: ir.DeploymentSpec{
				Definition: ir.DefinitionSpec{
					Metadata: ir.MetadataSpec{
						SpecVersion: ir.SpecVersion_0_2_0,
					},
				},
			},
			wantError: fmt.Errorf("invalid DAG, no sources found"),
		},
		{
			name: "too many sources",
			spec: ir.DeploymentSpec{
				Connectors: []ir.ConnectorSpec{
					{
						Type: ir.ConnectorSource,
					},
					{
						Type: ir.ConnectorSource,
					},
				},
			},
			wantError: fmt.Errorf("invalid DAG, too many sources"),
		},
		{
			name: "only one source",
			spec: ir.DeploymentSpec{
				Connectors: []ir.ConnectorSpec{
					{
						Type: ir.ConnectorSource,
					},
				},
			},
			wantError: fmt.Errorf("invalid DAG, there has to be at least one source and one destination"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spec := tc.spec
			_, err := spec.BuildDAG()
			if tc.wantError != nil {
				assert.Equal(t, err.Error(), tc.wantError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
