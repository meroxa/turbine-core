package ir

import (
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

func ValidateSpec(spec []byte, specVersion string) error {
	if specVersion != "0.1.1" {
		return fmt.Errorf("spec version %q is not supported", specVersion)
	}

	sch, err := jsonschema.CompileString("turbine.ir.schema.json", turbineIRSchema)
	if err != nil {
		return err
	}

	var v interface{}
	if err := json.Unmarshal(spec, &v); err != nil {
		return err
	}

	if err = sch.Validate(v); err != nil {
		return err
	}

	return nil
}

var turbineIRSchema = `
{
	"$schema": "https://json-schema.org/draft/2019-09/schema#",
	"$id": "https://meroxa.io/turbine.ir.schema.json",
	"title": "Turbine intermediate representation schema",
	"type": "object",
	"properties": {
		"connectors": {
			"type": "array",
			"items": [
				{
					"type": "object",
					"properties": {
						"collection": {
							"type": "string",
							"minLength": 1
						},
						"type": {
							"type": "string",
							"enum": ["source", "destination" ]
						},
						"resource": {
							"type": "string",
							"minLength": 1
						},
						"config": {
							"type": "object"
						}
					},
					"required": [
						"collection",
						"type",
						"resource"
					]
				}
			],
			"contains": {
				"type": "object",
				"properties": {
					"collection": {
						"type": "string"
					},
					"type": {
						"type": "string",
						"pattern": "^source$"
					},
					"resource": {
						"type": "string"
					},
					"config": {
						"type": "object"
					}
				},
				"required": ["type"]
			},
			"maxContains": 1,
			"minItems": 1,
			"uniqueItems": true
		},
		"functions": {
			"type": "array",
			"items": [
				{
					"type": "object",
					"properties": {
						"name": {
							"type": "string",
							"pattern": "^[a-zA-Z][a-zA-Z0-9-_]*$"
						},
						"image": {
							"type": "string",
							"minLength": 1
						},
						"env_vars": {
							"type": "object"
						}
					},
					"required": [ "name", "image" ]
				}
			],
			"minItems": 0,
			"maxItems": 1,
			"uniqueItems": true
		},
		"metadata": {
			"description": "The metadata associated with the spec",
			"type": "object",
			"properties": {
				"turbine": {
					"description": "The turbine details",
					"type": "object",
					"properties": {
						"language": {
							"description": "The language used to create deployment",
							"type": "string",
							"enum": ["golang", "js", "python" ]
						},
						"version": {
							"description": "The version of language used to create deployment",
							"type": "string",
							"minLength": 1
						}
					},
					"required": [ "language", "version" ]
				}
			},
			"required": [ "turbine" ]
		}
	},
  	"required": [
		"connectors",
		"metadata"
	]
}
`
