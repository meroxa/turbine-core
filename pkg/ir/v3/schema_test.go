package ir

import (
	"testing"

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
			specVersion: "v3",
			spec:        `{}`,
			err:         "\"\" field fails /required validation: missing properties: 'connectors', 'definition'",
		},
		{
			desc:        "empty definition",
			specVersion: "v3",
			spec: `{
						"connectors": [
							{
								"uuid": "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
								"name":   "my_source",
								"plugin_type": "source",
								"plugin_name": "postgres"
							}
						],
						"definition": {}
					}`,
			err: "\"/definition\" field fails /properties/definition/required validation: missing properties: 'git_sha', 'metadata'",
		},
		{
			desc:        "empty metadata",
			specVersion: "v3",
			spec: `{
						"connectors": [
							{
								"uuid": "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
								"name":   "my_source",
								"plugin_type": "source",
								"plugin_name": "postgres"
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
			specVersion: "v3",
			spec: `{
				"connectors": [
					{
						"uuid":   "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
						"name": "my-source",
						"collection": "users",
						"plugin_type": "source",
						"plugin_name": "postgres"
					}
				],
				"definition": {
					"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
					"metadata": {
						"turbine": {},
						"spec_version": "v3"
					}
				}
			}`,
			err: "\"/definition/metadata/turbine\" field fails /properties/definition/properties/metadata/properties/turbine/required validation: missing properties: 'language', 'version'",
		},
		{
			desc:        "minimal valid spec",
			specVersion: "v3",
			spec: `{
						"connectors": [
							{
								"name": "my_source",
								"plugin_type": "source",
								"plugin_name": "postgres"
							}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "v3"
							}
						}
					}`,
		},
		{
			desc:        "allow an empty connectors list",
			specVersion: "v3",
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
						"spec_version": "v3"
					}
				}
			}`,
		},
		{
			desc:        "empty connector",
			specVersion: "v3",
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
								"spec_version": "v3"
							}
						}
					}`,
			err: "\"/connectors/0\" field fails /properties/connectors/prefixItems/0/required validation: missing properties: 'name', 'plugin_type', 'plugin_name'",
		},
		{
			desc:        "unknown connector direction",
			specVersion: "v3",
			spec: `{
						"connectors": [
							{
								"uuid":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"name": "my_random_connector",
								"plugin_type": "some random direction",
								"plugin_name": "postgres"
							}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "v3"
							}
						}
					}`,
			err: "\"/connectors/0/plugin_type\" field fails /properties/connectors/prefixItems/0/properties/plugin_type/enum validation: value must be one of \"source\", \"destination\"",
		},
		{
			desc:        "allow one destination connector",
			specVersion: "v3",
			spec: `{
						"connectors": [
							{
								"uuid":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"name":   "my_destination",
								"plugin_type": "destination",
								"plugin_name": "postgres"
							}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "v3"
							}
						}
					}`,
		},
		{
			desc:        "one source, one destination connectors",
			specVersion: "v3",
			spec: `{
						"connectors": [
							{
								"name":   "my_destination",
								"plugin_type": "destination",
								"plugin_name": "postgres"
							},
							{
								"name":   "my_source",
								"plugin_type": "source",
								"plugin_name": "postgres"
							}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "v3"
							}
						}
					}`,
		},
		{
			desc:        "one source, two destination connectors",
			specVersion: "v3",
			spec: `{
						"connectors": [
							{
								"uuid":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"name":   "my_destination",
								"plugin_type": "destination",
								"plugin_name": "postgres"
							},
							{
								"uuid":   "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
								"name":   "my_destination",
								"plugin_type": "destination",
								"plugin_name": "postgres"
							},
							{
								"uuid":   "9e9e8e88-3a56-4a2a-993e-bfe49d526d07",
								"name":   "my_source",
								"plugin_type": "source",
								"plugin_name": "postgres"
							}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "v3"
							}
						}
					}`,
		},
		{
			desc:        "allow multiple sources, one destination connectors",
			specVersion: "v3",
			spec: `{
						"connectors": [
							{
								"uuid":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"name":   "my_destination",
								"plugin_type": "destination",
								"plugin_name": "postgres"
							},
							{
								"uuid":   "9839888cc-3a56-4a2a-993e-bfe49d526d07",
								"name":   "my_source",
								"plugin_type": "source",
								"plugin_name": "postgres"
							},
							{
								"uuid":   "02929383-3a56-4a2a-993e-bfe49d526d07",
								"name":   "my_source",
								"plugin_type": "source",
								"plugin_name": "postgres"
							}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "v3"
							}
						}
					}`,
		},
		{
			desc:        "one source, two duplicate destination connectors",
			specVersion: "v3",
			spec: `{
						"connectors": [
							{
								"uuid":   "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
								"plugin_type": "destination",
								"plugin_name": "postgres"
							},
							{
								"uuid":  "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
								"plugin_type": "destination",
								"plugin_name": "postgres"
							},
							{
								"uuid":  "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"plugin_type": "source",
								"plugin_name": "postgres"
							}
						],
						"definition": {
							"git_sha" : "83e7c39d83fe4cc04a404182dc30b8d9bed2537b",
							"metadata": {
								"turbine": {
									"language": "golang",
									"version": "0.19"
								},
								"spec_version": "v3"
							}
						}
					}`,
			err: "\"/connectors\" field fails /properties/connectors/uniqueItems validation: items at index 0 and 1 are equal",
		},
		{
			desc:        "empty function list",
			specVersion: "v3",
			spec: `{
						"connectors": [
							{
								"uuid":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"name":   "my_source",
								"plugin_type": "source",
								"plugin_name": "postgres"
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
								"spec_version": "v3"
							}
						}
					}`,
		},
		{
			desc:        "empty function",
			specVersion: "v3",
			spec: `{
						"connectors": [
							{
								"uuid":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"name":   "my_source",
								"plugin_type": "source",
								"plugin_name": "postgres"
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
								"spec_version": "v3"
							}
						}
					}`,
			err: "\"/functions/0\" field fails /properties/functions/prefixItems/0/required validation: missing properties: 'name', 'image'",
		},
		{
			desc:        "one function",
			specVersion: "v3",
			spec: `{
						"connectors": [
							{
								"uuid":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"name":   "my_source",
								"plugin_type": "source",
								"plugin_name": "postgres"
							}
						],
						"functions": [
							{
								"uuid": "d07f1a3d-f7e2-4495-a8fe-df46bef38a2b",
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
								"spec_version": "v3"
							}
						}
					}`,
		},
		{
			desc:        "two functions",
			specVersion: "v3",
			spec: `{
				"connectors": [
					{
						"uuid":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
						"name":   "my_source",
						"plugin_type": "source",
						"plugin_name": "postgres"
					}
				],
				"functions": [
					{
						"uuid": "d07f1a3d-f7e2-4495-a8fe-df46bef38a2b",
						"name": "enrich",
						"image": "ftorres/enrich:9"
					},
					{
						"uuid": "00d7f1a24-f7e2-4495-a8fe-df46bef32345",
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
						"spec_version": "v3"
					}
				}
			}`,
			err: "\"/functions\" field fails /properties/functions/maxItems validation: maximum 1 items required, but found 2 items",
		},
		{
			desc:        "maximum spec",
			specVersion: "v3",
			spec: `{
						"functions": [
							{
								"uuid": "d07f1a3d-f7e2-4495-a8fe-df46bef38a2b",
								"name": "anonymize",
								"image": "ec3b84a9-0866-4003-8e67-1492e9a3e61e"}
						],
						"connectors": [
							{
								"uuid":   "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
								"plugin_type": "source",
								"config": {},
								"plugin_name": "postgres",
								"name":   "my_source"
							},
							{
								"uuid":   "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
								"plugin_type": "destination",
								"config": {},
								"plugin_name": "postgres"
							}
						],
						"definition": {
							"git_sha": "b1537986d46bcd810960696d1e6df739e7bcc592",
							"metadata": {
								"turbine": {
									"version": "1.5.1",
									"language": "python"
								},
								"spec_version": "v3"
							}
						}
					}`,
		},
		{
			desc:        "spec with streams ",
			specVersion: "v3",
			spec: `{
					"secrets": {
						"API_KEY": "token"
					},
					"functions": [
						{
							"uuid": "d07f1a3d-f7e2-4495-a8fe-df46bef38a2b",
							"name": "anonymize",
							"image": "ec3b84a9-0866-4003-8e67-1492e9a3e61e"
						}
					],
					"connectors": [
						{
							"uuid": "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
							"plugin_type": "source",
							"config": {},
							"plugin_name": "postgres",
							"name":   "my_source"
						},
						{
							"uuid": "68dde1cc-3a56-4a2a-993e-bfe49d526d07",
							"plugin_type": "destination",
							"config": {},
							"plugin_name": "postgres",
							"name":   "my_destination"
						}
					],
					"stream": [
						{
							"uuid": "12345",
							"to_uuid": "d07f1a3d-f7e2-4495-a8fe-df46bef38a2b",
							"from_uuid": "13ae6f06-9fd0-4395-906e-9bba9a76ffc0",
							"name": "my_stream1"
						},
						{
							"uuid": "123456",
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
							"spec_version": "v3"
						}
					}
				}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := ValidateSpec([]byte(tc.spec), tc.specVersion)
			if tc.err == "" {
				require.NoError(t, err)
			} else {
				require.Equal(t, err.Error(), tc.err)
			}
		})
	}
}
