package main

import (
	"fmt"

	"github.com/meroxa/turbine-core/pkg/ir"
)


func main() {
	var spec ir.DeploymentSpec

	from := spec.AddSourceConnector(ir.ConnectorSpec{
		Type: ir.ConnectorSource,
		Resource: "source-resource",
		Collection: "events",
	})

	f1 := spec.AddFunction(ir.FunctionSpec{
		Name: "digest",
	}, from.ID)

	f2 := spec.AddFunction(ir.FunctionSpec{
		Name: "ingest",
	}, f1.ID)

	to1 := spec.AddDestinationConnector(ir.ConnectorSpec{
		Type: ir.ConnectorDestination,
		Resource: "ingested-resource",
		Collection: "ingested_events",
	}, f1.ID)

	to2 := spec.AddDestinationConnector(ir.ConnectorSpec{
		Type: ir.ConnectorDestination,
		Resource: "digested-resource",
		Collection: "digested_events",
	}, f2.ID)

	fmt.Printf("source: %+v\n", from)
	fmt.Printf("f1: %+v\n", f1)
	fmt.Printf("f2: %+v\n", f2)
	fmt.Printf("to2: %+v\n", to1)
	fmt.Printf("to2: %+v\n", to2)

	Walk(spec)
}

func Walk(spec ir.DeploymentSpec) {
	// find source
	// find stream of source
	// walk streams -> things

	var root ir.ConnectorSpec
	for _, c := range spec.Connectors {
		if c.Type == ir.ConnectorSource {
			root = c
			break
		}
	}

	fmt.Printf("root: %+v\n", root)

	for _, s := range spec.Streams {
		if s.FromID == root.ID {
			fmt.Printf("root %+v -> %+v\n", s.FromID, s.ToID)
		}
	}
}