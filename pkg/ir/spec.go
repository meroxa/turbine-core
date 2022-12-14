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
	mu          sync.Mutex
	turbineDag  *dag.DAG
	dagInitOnce sync.Once
	Secrets     map[string]string `json:"secrets,omitempty"`
	Connectors  []ConnectorSpec   `json:"connectors"`
	Functions   []FunctionSpec    `json:"functions,omitempty"`
	Streams     []StreamSpec      `json:"streams,omitempty"`
	Definition  DefinitionSpec    `json:"definition"`
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
	d.init()

	if len(d.turbineDag.GetRoots()) >= 1 {
		return fmt.Errorf("source connector already exists, can only add one per application")
	}
	if c.Type != ConnectorSource {
		return fmt.Errorf("connector type isn't a source, please check you are reading from a source connector.")
	}
	d.Connectors = append(d.Connectors, *c)
	return d.turbineDag.AddVertexByID(c.UUID, &c)
}

func (d *DeploymentSpec) AddFunction(f *FunctionSpec) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.init()

	d.Functions = append(d.Functions, *f)
	err := d.turbineDag.AddVertexByID(f.UUID, &f)
	if err != nil {
		return err
	}
	return nil
}

func (d *DeploymentSpec) AddDestination(c *ConnectorSpec) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.init()

	if c.Type != ConnectorDestination {
		return fmt.Errorf("connector type isn't a destination, please check you are writing to destination connector.")
	}
	d.Connectors = append(d.Connectors, *c)
	return d.turbineDag.AddVertexByID(c.UUID, &c)
}

func (d *DeploymentSpec) AddStream(s *StreamSpec) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.init()

	if _, err := d.turbineDag.GetVertex(s.FromUUID); err != nil {
		return fmt.Errorf("source with UUID - (%s) does not exist", s.FromUUID)
	}

	if _, err := d.turbineDag.GetVertex(s.ToUUID); err != nil {
		return fmt.Errorf("destination with UUID - (%s) does not exist", s.ToUUID)
	}

	d.Streams = append(d.Streams, *s)
	if err := d.turbineDag.AddEdge(s.FromUUID, s.ToUUID); err != nil {
		return err
	}
	return nil
}

func (d *DeploymentSpec) ValidateDag() error {
	if len(d.turbineDag.GetRoots()) > 1 {
		return fmt.Errorf("more than one source / root detected, can only add one per application. please ensure your resources are connected.")
	}
	return nil
}

func (d *DeploymentSpec) BuildDAG() string {
	return d.turbineDag.String()
}

func (d *DeploymentSpec) init() {
	d.dagInitOnce.Do(func() {
		d.turbineDag = dag.NewDAG()
	})
}
