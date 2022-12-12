package ir

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type ConnectorType string
type Lang string

const (
	GoLang               Lang          = "golang"
	JavaScript           Lang          = "javascript"
	NodeJs               Lang          = "nodejs"
	Python               Lang          = "python"
	Python3              Lang          = "python3"
	Ruby                 Lang          = "ruby"
	ConnectorSource      ConnectorType = "source"
	ConnectorDestination ConnectorType = "destination"

	LatestSpecVersion = "0.2.0"
)

type DeploymentSpec struct {
	mu            sync.Mutex
	DeploymentMap DeploymentMap
	Secrets       map[string]string `json:"secrets,omitempty"`
	Connectors    []ConnectorSpec   `json:"connectors"`
	Functions     []FunctionSpec    `json:"functions,omitempty"`
	Streams       []StreamSpec      `json:"streams,omitempty"`
	Definition    DefinitionSpec    `json:"definition"`
}

type DeploymentMap struct {
	Source *ConnectorSpec
	Nodes  []Node
}
type Node struct {
	UUID  string
	Edges map[string]string
}

type StreamSpec struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	FromUUID string `json:"from_uuid"`
	ToUUID   string `json:"to_uuid"`
}

type ConnectorSpec struct {
	UUID       string                 `json:"uuid"`
	Type       ConnectorType          `json:"type"`
	Resource   string                 `json:"resource"`
	Collection string                 `json:"collection"`
	Config     map[string]interface{} `json:"config,omitempty"`
}

type FunctionSpec struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type DefinitionSpec struct {
	GitSha   string       `json:"git_sha"`
	Metadata MetadataSpec `json:"metadata"`
}

type MetadataSpec struct {
	Turbine     TurbineSpec `json:"turbine"`
	SpecVersion string      `json:"spec_version"`
}

type TurbineSpec struct {
	Language Lang   `json:"language"`
	Version  string `json:"version"`
}

func ValidateSpecVersion(specVersion string) error {
	if specVersion != LatestSpecVersion {
		return fmt.Errorf("spec version %q is not a supported. use version %q instead", specVersion, LatestSpecVersion)
	}
	return nil
}

func (s *DeploymentSpec) ValidateStream() error {
	for _, stream := range s.Streams {
		if stream.FromUUID == stream.ToUUID {
			return fmt.Errorf("for stream %q , ids for source (%q) and destination (%q) must be different.", stream.Name, stream.FromUUID, stream.ToUUID)
		}
	}
	return nil
}

func (s *DeploymentSpec) SetImageForFunctions(image string) {
	for i := range s.Functions {
		s.Functions[i].Image = image
	}
}

func (s *DeploymentSpec) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func Unmarshal(data []byte) (*DeploymentSpec, error) {
	spec := &DeploymentSpec{}
	if err := json.Unmarshal(data, spec); err != nil {
		return nil, err
	}
	return spec, nil
}

func (d *DeploymentSpec) AddSource(c *ConnectorSpec) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// check if spec is source
	if d.DeploymentMap.Source != nil {
		return fmt.Errorf("source connector (%s) already exists", d.DeploymentMap.Source.UUID)
	}

	for _, node := range d.DeploymentMap.Nodes {
		if node.Item.UUID == c.UUID {
			return fmt.Errorf("connector (%s) already added", c.UUID)
		}
	}

	d.Connectors = append(d.Connectors, *c)
	d.DeploymentMap.Nodes = append(d.DeploymentMap.Nodes, Node{
		Item:  c,
		Edges: make(map[string]string),
	})
	//track source connector spec
	d.DeploymentMap.Source = c

	return nil
}

func (d *DeploymentSpec) AddDestination(c *ConnectorSpec) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i, node := range d.DeploymentMap.Nodes {
		if node == c.UUID {
			return fmt.Errorf("connector (%s) already added", c.UUID)
		}
	}

	d.Connectors = append(d.Connectors, *c)
	d.DeploymentMap.Nodes = append(d.DeploymentMap.Nodes, Node{
		UUID:  c.UUID,
		Edges: make(map[string]interface{}),
	})
	return nil
}

// .. similarly for functions

func (d *DeploymentSpec) AddStream(fromUUID, toUUID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// check if stream exists based on ID
	// check if from id / to id are valid

	for _, node := range d.DeploymentMap.Nodes {
		if node.UUID == fromUUID {
			return fmt.Errorf("connector (%s) already added", c.UUID)
		}
	}

	d.Streams = append(d.Streams, StreamSpec{
		UUID:     uuid.New(),
		FromUUID: fromUUID,
		ToUUID:   toUUID,
		Name:     fromUUID + "_" + toUUID,
	})
	// TODO: cache edges
	return nil
}

func (d *DeploymentSpec) BuildDAG() error {

	//add all functions as nodes
	for _, f := range d.Functions {
		d.DeploymentMap.Nodes = append(d.DeploymentMap.Nodes, Node{
			UUID:  f.UUID,
			Edges: make(map[string]interface{}),
		})
	}

	//add all streams as nodes
	for _, s := range d.Streams {
		d.DeploymentMap.Nodes = append(d.DeploymentMap.Nodes, Node{
			UUID:  s.UUID,
			Edges: make(map[string]interface{}),
		})
	}
	//connectors are already added so no need

	// walk all connectors and add them to the d.nodes
	// ensure single source
	// walk all functions and add them to the d.nodes
	// walk all streams and add them to d.nodes
	// return error on:
	// * more than one source
	// * ID collision for streams/functions/connectors
	// * stream FromID/ToID

	if d.DeploymentMap.Source == nil {

	}

	//if err := spec.BuildDAG(); err != nil {
	//}
	// failed to build the dag
	return nil
}

func (d *DeploymentSpec) AddEdge(node1, node2 Node) error {
	d.DeploymentMap.Nodes[node1.UUID].Edges[node2.UUID]
	g.nodes[n1].edges[n2] = w
	return nil
}
