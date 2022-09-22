package ir

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

var (
	//go:embed schema.json
	turbineIRSchema string
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
