package ir

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/heimdalr/dag"
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
	mu         sync.Mutex
	turbineDAG turbineDAG
	Secrets    map[string]string `json:"secrets,omitempty"`
	Connectors []ConnectorSpec   `json:"connectors"`
	Functions  []FunctionSpec    `json:"functions,omitempty"`
	Streams    []StreamSpec      `json:"streams,omitempty"`
	Definition DefinitionSpec    `json:"definition"`
}

type turbineDAG struct {
	dag    *dag.DAG
	source *ConnectorSpec
	nodes  map[string]string
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

type idInterface interface {
	ID() string
}

func ValidateSpecVersion(specVersion string) error {
	if specVersion != LatestSpecVersion {
		return fmt.Errorf("spec version %q is not a supported. use version %q instead", specVersion, LatestSpecVersion)
	}
	return nil
}

func (d *DeploymentSpec) SetImageForFunctions(image string) {
	for i := range d.Functions {
		d.Functions[i].Image = image
	}
}

func (d *DeploymentSpec) Marshal() ([]byte, error) {
	return json.Marshal(d)
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
	if d.turbineDAG.source != nil {
		return fmt.Errorf("source connector (%s) already exists", d.turbineDAG.source.UUID)
	}

	if source := d.turbineDAG.nodes[c.UUID]; source != "" {
		return fmt.Errorf("connector with uuid %s already exists", c.UUID)
	}
	d.Connectors = append(d.Connectors, *c)
	d.turbineDAG.source = c
	src, err := d.turbineDAG.dag.AddVertex(&c)
	if err != nil {
		return err
	}
	d.turbineDAG.nodes[c.UUID] = src
	return nil
}

func (d *DeploymentSpec) AddFunction(f *FunctionSpec) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if function := d.turbineDAG.nodes[f.UUID]; function != "" {
		return fmt.Errorf("function with uuid %s already exists", f.UUID)
	}
	d.Functions = append(d.Functions, *f)
	fun, err := d.turbineDAG.dag.AddVertex(&f)
	if err != nil {
		return err
	}
	d.turbineDAG.nodes[f.UUID] = fun
	return nil
}

func (d *DeploymentSpec) AddDestination(c *ConnectorSpec) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Connectors = append(d.Connectors, *c)
	if dest := d.turbineDAG.nodes[c.UUID]; dest != "" {
		return fmt.Errorf("connector with uuid %s already exists", c.UUID)
	}
	dest, err := d.turbineDAG.dag.AddVertex(&c)
	if err != nil {
		return err
	}
	d.turbineDAG.nodes[c.UUID] = dest
	return nil
}

func (d *DeploymentSpec) AddStream(s *StreamSpec) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	source, ok := d.turbineDAG.nodes[s.FromUUID]
	if !ok {
		return fmt.Errorf("source node (%s) does not exist", s.FromUUID)
	}

	dest, ok := d.turbineDAG.nodes[s.ToUUID]

	if !ok {
		return fmt.Errorf("destination node (%s) does not exist", s.ToUUID)
	}

	d.Streams = append(d.Streams, *s)
	if err := d.turbineDAG.dag.AddEdge(source, dest); err != nil {
		return err
	}
	return nil
}

func (d *DeploymentSpec) InitDag() {
	d.turbineDAG = turbineDAG{
		nodes: make(map[string]string),
		dag:   dag.NewDAG(),
	}
}

func (f FunctionSpec) ID() string {
	return f.UUID
}
func (s StreamSpec) ID() string {
	return s.UUID
}
func (c ConnectorSpec) ID() string {
	return c.UUID
}
