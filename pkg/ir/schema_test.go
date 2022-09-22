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
			specVersion: "0.1.1",
			spec:        `{}`,
			err:         "missing properties: 'connectors', 'metadata'",
		},
		{
			desc:        "empty metadata",
			specVersion: "0.1.1",
			spec: `{
				"connectors": [
					{
						"collection": "users",
						"type": "source",
						"resource": "pg"
					}
				],
				"metadata": {}
			}`,
			err: "missing properties: 'turbine'",
		},
		{
			desc:        "empty turbine",
			specVersion: "0.1.1",
			spec: `{
				"connectors": [
					{
						"collection": "users",
						"type": "source",
						"resource": "pg"
					}
				],
				"metadata": {
					"turbine": {}
				}
			}`,
			err: "missing properties: 'language', 'version'",
		},
		{
			desc:        "minimal valid spec",
			specVersion: "0.1.1",
			spec: `{
				"connectors": [
					{
						"collection": "users",
						"type": "source",
						"resource": "pg"
					}
				],
				"metadata": {
					"turbine": {
						"language": "golang",
						"version": "0.19"
					}
				}
			}`,
		},
		{
			desc:        "connectors list",
			specVersion: "0.1.1",
			spec: `{
				"connectors": [
				],
				"metadata": {
					"turbine": {
						"language": "golang",
						"version": "0.19"
					}
				}
			}`,
			err: "minimum 1 items required, but found 0 items",
		},
		{
			desc:        "empty connector",
			specVersion: "0.1.1",
			spec: `{
				"connectors": [
					{}
				],
				"metadata": {
					"turbine": {
						"language": "golang",
						"version": "0.19"
					}
				}
			}`,
			err: "missing properties: 'collection', 'type', 'resource'",
		},
		{
			desc:        "unknown connector type",
			specVersion: "0.1.1",
			spec: `{
				"connectors": [
					{
						"collection": "users",
						"type": "some random type",
						"resource": "pg"
					}
				],
				"metadata": {
					"turbine": {
						"language": "golang",
						"version": "0.19"
					}
				}
			}`,
			err: "value must be one of \"source\", \"destination\"",
		},
		{
			desc:        "one destination connector",
			specVersion: "0.1.1",
			spec: `{
				"connectors": [
					{
						"collection": "users",
						"type": "destination",
						"resource": "pg"
					}
				],
				"metadata": {
					"turbine": {
						"language": "golang",
						"version": "0.19"
					}
				}
			}`,
			err: "does not match pattern '^source$'",
		},
		{
			desc:        "one source, one destination connectors",
			specVersion: "0.1.1",
			spec: `{
				"connectors": [
					{
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
				"metadata": {
					"turbine": {
						"language": "golang",
						"version": "0.19"
					}
				}
			}`,
		},
		{
			desc:        "one source, two destination connectors",
			specVersion: "0.1.1",
			spec: `{
				"connectors": [
					{
						"collection": "users_processed",
						"type": "destination",
						"resource": "pg"
					},
					{
						"collection": "users_copy",
						"type": "destination",
						"resource": "pg"
					},
					{
						"collection": "users",
						"type": "source",
						"resource": "pg"
					}
				],
				"metadata": {
					"turbine": {
						"language": "golang",
						"version": "0.19"
					}
				}
			}`,
		},
		{
			desc:        "two source, one destination connectors",
			specVersion: "0.1.1",
			spec: `{
				"connectors": [
					{
						"collection": "users_processed",
						"type": "destination",
						"resource": "pg"
					},
					{
						"collection": "accounts",
						"type": "source",
						"resource": "pg"
					},
					{
						"collection": "users",
						"type": "source",
						"resource": "pg"
					}
				],
				"metadata": {
					"turbine": {
						"language": "golang",
						"version": "0.19"
					}
				}
			}`,
			err: "valid must be <= 1, but got 2",
		},
		{
			desc:        "one source, two duplicate destination connectors",
			specVersion: "0.1.1",
			spec: `{
				"connectors": [
					{
						"collection": "users_processed",
						"type": "destination",
						"resource": "pg"
					},
					{
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
				"metadata": {
					"turbine": {
						"language": "golang",
						"version": "0.19"
					}
				}
			}`,
			err: "items at index 0 and 1 are equal",
		},
		{
			desc:        "empty function list",
			specVersion: "0.1.1",
			spec: `{
				"connectors": [
					{
						"collection": "users",
						"type": "source",
						"resource": "pg"
					}
				],
				"function": [],
				"metadata": {
					"turbine": {
						"language": "golang",
						"version": "0.19"
					}
				}
			}`,
		},
		{
			desc:        "empty function",
			specVersion: "0.1.1",
			spec: `{
				"connectors": [
					{
						"collection": "users",
						"type": "source",
						"resource": "pg"
					}
				],
				"functions": [
					{}
				],
				"metadata": {
					"turbine": {
						"language": "golang",
						"version": "0.19"
					}
				}
			}`,
			err: "missing properties: 'name', 'image'",
		},
		{
			desc:        "one function",
			specVersion: "0.1.1",
			spec: `{
				"connectors": [
					{
						"collection": "users",
						"type": "source",
						"resource": "pg"
					}
				],
				"functions": [
					{
						"name": "enrich",
						"image": "ftorres/enrich:9"
					}
				],
				"metadata": {
					"turbine": {
						"language": "golang",
						"version": "0.19"
					}
				}
			}`,
		},
		{
			desc:        "two functions",
			specVersion: "0.1.1",
			spec: `{
				"connectors": [
					{
						"collection": "users",
						"type": "source",
						"resource": "pg"
					}
				],
				"functions": [
					{
						"name": "enrich",
						"image": "ftorres/enrich:9"
					},
					{
						"name": "enrich_new",
						"image": "ftorres/enrich:1000"
					}
				],
				"metadata": {
					"turbine": {
						"language": "golang",
						"version": "0.19"
					}
				}
			}`,
			err: "maximum 1 items required, but found 2 items",
		},
		{
			desc:        "maximum spec",
			specVersion: "0.1.1",
			spec: `{
				"connectors": [
					{
						"collection": "users",
						"type": "source",
						"resource": "pg",
						"config": {
							"logical_replication": true
						}
					}, 
					{
						"collection": "users_enriched",
						"type": "destination",
						"resource": "pg_prod"
					}
				],
				"functions": [
					{
						"name": "enrich",
						"image": "ftorres/enrich:9",
						"env_vars": {
							"CLEARBIT_API_KEY": "token-1"
						}
					}
				],
				"metadata": {
					"turbine": {
						"language": "golang",
						"version": "0.19"
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
				require.Contains(t, err.Error(), tc.err)
			}
		})
	}
}
