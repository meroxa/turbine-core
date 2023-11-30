package ir

import (
	_ "embed"
	"encoding/json"
	ir "github.com/meroxa/turbine-core/pkg/ir/v2"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed schema.json
var turbineIRSchema string

func ValidateSpec(spec []byte, specVersion string) error {
	err := ir.ValidateSpecVersion(specVersion)
	if err != nil {
		return err
	}

	sch, err := jsonschema.CompileString("turbine.ir.schema.json", turbineIRSchema)
	if err != nil {
		return err
	}

	var v interface{}
	if err := json.Unmarshal(spec, &v); err != nil {
		return err
	}

	if err := sch.Validate(v); err != nil {
		return ir.NewSpecValidationError(err)
	}

	return nil
}
