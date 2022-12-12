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

func Test_EnsureSingleSource(t *testing.T) {
	var spec ir.DeploymentSpec
	spec.DeploymentMap = *spec.InitDag()
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
			UUID:       "2",
			Collection: "accounts2",
			Resource:   "mongo",
			Type:       ir.ConnectorSource,
			Config: map[string]interface{}{
				"config": "value",
			},
		},
	)

	gotError := spec.BuildDAG()
	fmt.Println(gotError)

	require.Error(t, gotError)

}

// Scenario 1 - Simple DAG
// ( src_con ) → (stream) → (function) → (stream) → (dest1)

func Test_Scenario1(t *testing.T) {
	var spec ir.DeploymentSpec
	spec.DeploymentMap = *spec.InitDag()
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

	gotError := spec.BuildDAG()

	fmt.Println(gotError)

}

// Scenario 2 - DAG with two Destinations from 1 function
//    ( src_con ) → (stream) → (function) → (stream) → (dest1)
//									↓
// 								(stream) → (dest2)

func Test_DAGScenario2(t *testing.T) {
	var spec ir.DeploymentSpec
	spec.DeploymentMap = *spec.InitDag()
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

	gotError := spec.BuildDAG()

	require.NoError(t, gotError)

}

// Scenario 3 - Not acyclic, trying to write from one destination to another?
//    ( src_con ) → (stream) → (function) → (stream) →  → (dest1)
//									↓
// 						    	(stream) →	(dest2) → (stream) ↑

func Test_DAGScenario3(t *testing.T) {
	var spec ir.DeploymentSpec
	spec.DeploymentMap = *spec.InitDag()
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

	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "4_3",
			Name:     "my_stream5",
			FromUUID: "4",
			ToUUID:   "3",
		},
	)

	require.NoError(t, err)

	gotError := spec.BuildDAG()

	require.Error(t, gotError)

}

// Scenario 4 - Not acyclic, trying to write from one function back to source
//    ( src_con ) → (stream) → (function 1) → (stream) → (function2)  → (dest1)
//			↑					    ↓
// 		    ←  ←  ← (stream) ← ← ← ←

func Test_DAGScenario4(t *testing.T) {
	var spec ir.DeploymentSpec
	spec.DeploymentMap = *spec.InitDag()
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

	gotError := spec.BuildDAG()

	require.Error(t, gotError)

}

// Scenario 5 - DAG, multuple functions, 1 destination
//    ( src_con ) → (stream) → (function 1) → (stream) → (function2)  → (dest1)

func Test_DAGScenario5(t *testing.T) {
	var spec ir.DeploymentSpec
	spec.DeploymentMap = *spec.InitDag()
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

	gotError := spec.BuildDAG()

	require.NoError(t, gotError)

}

// Scenario 6 - DAG with many functions and destinations
//    ( src_con ) → (stream) → (function 1) → (stream) → (function2)  → (dest1)
//								    ↓
//								 (stream)  → (dest2)
//								    ↓
//								(function 3)
//									↓
//								(stream)
//   								↓
//								(dest 3)

func Test_DAGScenario6(t *testing.T) {
	var spec ir.DeploymentSpec
	spec.DeploymentMap = *spec.InitDag()
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
	err = spec.AddStream(
		&ir.StreamSpec{
			UUID:     "6_7",
			Name:     "my_stream6",
			FromUUID: "6",
			ToUUID:   "7",
		},
	)

	gotError := spec.BuildDAG()

	require.NoError(t, gotError)

}
