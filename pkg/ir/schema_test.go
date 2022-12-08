package ir_test

import (
	"testing"

	"github.com/meroxa/turbine-core/pkg/ir"
	"github.com/stretchr/testify/require"
)

func Test_ValidSpec(t *testing.T) {
	testCases := []struct {
		desc        string
		specVersion string
		spec        string
		err         string
	}{
		{
			desc:        "empty spec",
			specVersion: "0.2.0",
			spec:        `{}`,
			err:         "\"\" field fails /required validation: missing properties: 'connectors', 'definition'",
		},
		{
			desc:        "empty definition",
			specVersion: "0.2.0",
			spec: `{
						"connectors": [
							{
								"id": "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
								"collection": "users",
								"type": "source",
								"resource": "pg"
							}
						],
						"definition": {}
					}`,
			err: "\"/definition\" field fails /properties/definition/required validation: missing properties: 'git_sha', 'metadata'",
		},
		{
			desc:        "empty metadata",
			specVersion: "0.2.0",
			spec: `{
						"connectors": [
							{
								"id": "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
								"collection": "users",
								"type": "source",
								"resource": "pg"
							}
						],
						"definition": {
							"git_sha": "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {}
						}
					}`,
			err: "\"/definition/metadata\" field fails /properties/definition/properties/metadata/required validation: missing properties: 'turbine', 'spec_version'",
		},
		{
			desc:        "empty turbine",
			specVersion: "0.2.0",
			spec: `{
				"connectors": [
					{
						"id":   "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
						"collection": "users",
						"type": "source",
						"resource": "pg"
					}
				],
				"definition": {
					"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
					"metadata": {
						"turbine": {},
						"spec_version": "0.2.0"
					}
				}
			}`,
			err: "\"/definition/metadata/turbine\" field fails /properties/definition/properties/metadata/properties/turbine/required validation: missing properties: 'language', 'version'",
		},
		{
			desc:        "minimal valid spec",
			specVersion: "0.2.0",
			spec: `{
						"connectors": [
							{
								"id":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"collection": "users",
								"type": "source",
								"resource": "pg"
							}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "0.2.0"
							}
						}
					}`,
		},
		{
			desc:        "connectors list",
			specVersion: "0.2.0",
			spec: `{
				"connectors": [
				],
				"definition": {
					"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
					"metadata": {
						"turbine": {
							"language": "golang",
							"version": "0.19"
						},
						"spec_version": "0.2.0"
					}
				}
			}`,
			err: "\"/connectors\" field fails /properties/connectors/minItems validation: minimum 1 items required, but found 0 items",
		},
		{
			desc:        "empty connector",
			specVersion: "0.2.0",
			spec: `{
						"connectors": [
							{}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "0.2.0"
							}
						}
					}`,
			err: "\"/connectors/0\" field fails /properties/connectors/items/0/required validation: missing properties: 'collection', 'type', 'resource'",
		},
		{
			desc:        "unknown connector type",
			specVersion: "0.2.0",
			spec: `{
						"connectors": [
							{
								"id":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"collection": "users",
								"type": "some random type",
								"resource": "pg"
							}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "0.2.0"
							}
						}
					}`,
			err: "\"/connectors/0/type\" field fails /properties/connectors/items/0/properties/type/enum validation: value must be one of \"source\", \"destination\"",
		},
		{
			desc:        "one destination connector",
			specVersion: "0.2.0",
			spec: `{
						"connectors": [
							{
								"id":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"collection": "users",
								"type": "destination",
								"resource": "pg"
							}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "0.2.0"
							}
						}
					}`,
			err: "\"/connectors/0/type\" field fails /properties/connectors/contains/properties/type/pattern validation: does not match pattern '^source$'",
		},
		{
			desc:        "one source, one destination connectors",
			specVersion: "0.2.0",
			spec: `{
						"connectors": [
							{
								"id":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"collection": "users_processed",
								"type": "destination",
								"resource": "pg"
							},
							{
								"collection": "users",
								"type": "source",
								"resource": "pg"
							}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "0.2.0"
							}
						}
					}`,
		},
		{
			desc:        "one source, two destination connectors",
			specVersion: "0.2.0",
			spec: `{
						"connectors": [
							{
								"id":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"collection": "users_processed",
								"type": "destination",
								"resource": "pg"
							},
							{
								"id":   "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
								"collection": "users_copy",
								"type": "destination",
								"resource": "pg"
							},
							{
								"id":   "9e9e8e88-3a56-4a2a-993e-bfe49d526d07",
								"collection": "users",
								"type": "source",
								"resource": "pg"
							}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "0.2.0"
							}
						}
					}`,
		},
		{
			desc:        "two source, one destination connectors",
			specVersion: "0.2.0",
			spec: `{
						"connectors": [
							{
								"id":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"collection": "users_processed",
								"type": "destination",
								"resource": "pg"
							},
							{
								"id":   "9839888cc-3a56-4a2a-993e-bfe49d526d07",
								"collection": "accounts",
								"type": "source",
								"resource": "pg"
							},
							{
								"id":   "02929383-3a56-4a2a-993e-bfe49d526d07",
								"collection": "users",
								"type": "source",
								"resource": "pg"
							}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "0.2.0"
							}
						}
					}`,
			err: "\"/connectors\" field fails /properties/connectors/maxContains validation: valid must be <= 1, but got 2",
		},
		{
			desc:        "one source, two duplicate destination connectors",
			specVersion: "0.2.0",
			spec: `{
						"connectors": [
							{
								"id":   "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
								"collection": "users_processed",
								"type": "destination",
								"resource": "pg"
							},
							{
								"id":  "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
								"collection": "users_processed",
								"type": "destination",
								"resource": "pg"
							},
							{
								"id":  "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"collection": "users",
								"type": "source",
								"resource": "pg"
							}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "0.2.0"
							}
						}
					}`,
			err: "\"/connectors\" field fails /properties/connectors/uniqueItems validation: items at index 0 and 1 are equal",
		},
		{
			desc:        "empty function list",
			specVersion: "0.2.0",
			spec: `{
						"connectors": [
							{
								"id":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"collection": "users",
								"type": "source",
								"resource": "pg"
							}
						],
						"function": [],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "0.2.0"
							}
						}
					}`,
		},
		{
			desc:        "empty function",
			specVersion: "0.2.0",
			spec: `{
						"connectors": [
							{
								"id":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"collection": "users",
								"type": "source",
								"resource": "pg"
							}
						],
						"functions": [
							{}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "0.2.0"
							}
						}
					}`,
			err: "\"/functions/0\" field fails /properties/functions/items/0/required validation: missing properties: 'name', 'image'",
		},
		{
			desc:        "one function",
			specVersion: "0.2.0",
			spec: `{
						"connectors": [
							{
								"id":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"collection": "users",
								"type": "source",
								"resource": "pg"
							}
						],
						"functions": [
							{
								"id": "d07f1a3d-f7e2-4495-a8fe-df46bef38a2b",
								"name": "enrich",
								"image": "ftorres/enrich:9"
							}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "0.2.0"
							}
						}
					}`,
		},
		{
			desc:        "two functions",
			specVersion: "0.2.0",
			spec: `{
						"connectors": [
							{
								"id":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"collection": "users",
								"type": "source",
								"resource": "pg"
							}
						],
						"functions": [
							{
								"id": "d07f1a3d-f7e2-4495-a8fe-df46bef38a2b",
								"name": "enrich",
								"image": "ftorres/enrich:9"
							},
							{
								"id": "00d7f1a24-f7e2-4495-a8fe-df46bef32345",
								"name": "enrich_new",
								"image": "ftorres/enrich:1000"
							}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "0.2.0"
							}
						}
					}`,
		},
		{
			desc:        "maximum spec",
			specVersion: "0.2.0",
			spec: `{
						"secrets": {
							"API_KEY": "token"
						}, 
						"functions": [
							{
								"id": "d07f1a3d-f7e2-4495-a8fe-df46bef38a2b",
								"name": "anonymize",
								"image": "ec3b84a9-0866-4003-8e67-1492e9a3e61e"}
						],
						"connectors": [
							{
								"id":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"type": "source",
								"config": {}, 
								"resource": "pg", 
								"collection": "sequences"
							}, 
							{
								"id":   "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
								"type": "destination", 
								"config": {}, 
								"resource": "pg", 
								"collection": "test_py_feature_branch"
							}
						], 
						"definition": {
							"git_sha": "b1537986d46bcd810960696d1e6df739e7bcc592", 
							"metadata": {
								"turbine": {
									"version": "1.5.1", 
									"language": "py"
								}, 
								"spec_version": "0.2.0"
							}
						}
					}`,
		},
		{
			desc:        "spec with streams ",
			specVersion: "0.2.0",
			spec: `{
					"secrets": {
						"API_KEY": "token"
					},
					"functions": [
						{
							"id": "d07f1a3d-f7e2-4495-a8fe-df46bef38a2b",
							"name": "anonymize",
							"image": "ec3b84a9-0866-4003-8e67-1492e9a3e61e"
						}
					],
					"connectors": [
						{
							"id": "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
							"type": "source",
							"config": {},
							"resource": "pg",
							"collection": "sequences"
						},
						{
							"id": "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
							"type": "destination",
							"config": {},
							"resource": "pg",
							"collection": "test_py_feature_branch"
						}
					],
					"stream": [
						{
							"id": "12345",
							"to_uuid": "d07f1a3d-f7e2-4495-a8fe-df46bef38a2b",
							"from_uuid": "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
							"name": "my_stream1"
						},
						{
							"id": "123456",
							"from_uuid": "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
							"to_uuid": "d07f1a3d-f7e2-4495-a8fe-df46bef38a2b",
							"name": "my_stream2"
						}
					],
					"definition": {
						"git_sha": "b1537986d46bcd810960696d1e6df739e7bcc592",
						"metadata": {
							"turbine": {
								"version": "1.5.1",
								"language": "py"
							},
							"spec_version": "0.2.0"
						}
					}
				}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := ir.ValidateSpec([]byte(tc.spec), tc.specVersion)
			if tc.err == "" {
				require.NoError(t, err)
			} else {
				require.Equal(t, err.Error(), tc.err)
			}
		})
	}
}
