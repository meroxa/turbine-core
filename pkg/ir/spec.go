package ir

import (
	"encoding/json"
	"fmt"
	"sync"
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
	Nodes  map[string]interface{}
	Edges  map[string]map[string]interface{}
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

type IDInterface interface {
	ID() string
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
		if id, ok := node.(IDInterface); ok {
			if id.ID() == c.UUID {
				return fmt.Errorf("connector (%s) already added", c.UUID)
			}
		}
	}

	d.Connectors = append(d.Connectors, *c)
	d.DeploymentMap.Nodes[c.UUID] = c
	return nil
}

func (d *DeploymentSpec) AddFunction(f *FunctionSpec) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, node := range d.DeploymentMap.Nodes {
		if id, ok := node.(IDInterface); ok {
			if id.ID() == f.UUID {
				return fmt.Errorf("function (%s) already added", f.UUID)
			}
		}
	}

	d.Functions = append(d.Functions, *f)
	d.DeploymentMap.Nodes[f.UUID] = f
	return nil
}

func (d *DeploymentSpec) AddDestination(c *ConnectorSpec) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, node := range d.DeploymentMap.Nodes {
		if id, ok := node.(IDInterface); ok {
			if id.ID() == c.UUID {
				return fmt.Errorf("connector (%s) already added", c.UUID)
			}
		}
	}
	d.Connectors = append(d.Connectors, *c)
	d.DeploymentMap.Nodes[c.UUID] = c
	return nil
}

// .. similarly for functions

func (d *DeploymentSpec) AddStream(s *StreamSpec) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.DeploymentMap.Nodes[s.FromUUID]; !ok {
		return fmt.Errorf("source node (%s) does not exist", s.FromUUID)
	}

	if _, ok := d.DeploymentMap.Nodes[s.ToUUID]; !ok {
		return fmt.Errorf("destination node (%s) does not exist", s.ToUUID)
	}

	if err := d.ValidateStream(); err != nil {
		return err
	}

	d.DeploymentMap.Nodes[s.UUID] = s
	d.Streams = append(d.Streams, *s)

	if err := d.AddEdge(s.FromUUID, s.UUID, *s); err != nil {
		return err
	}
	if err := d.AddEdge(s.UUID, s.ToUUID, *s); err != nil {
		return err
	}

	return nil
}

func (d *DeploymentSpec) InitDag() *DeploymentMap {
	return &DeploymentMap{
		Nodes: make(map[string]interface{}),
		Edges: make(map[string]map[string]interface{}),
		//Source: d.DeploymentMap.Source,
	}
}

func (d *DeploymentSpec) BuildDAG() error {
	// add all connectors as nodes
	// also get and set source

	var source ConnectorSpec
	for _, c := range d.Connectors {
		d.DeploymentMap.Nodes[c.UUID] = c
		if c.Type == ConnectorSource {
			if source.UUID == "" {
				source = c
			} else {
				return fmt.Errorf("found more than one source, cannot add source with UUID %s as source with UUID %s already exists", c.UUID, source.UUID)
			}
		}
	}

	d.DeploymentMap.Source = &source

	fmt.Println("Source:")
	fmt.Println(d.DeploymentMap.Source)

	//add all functions as nodes
	for _, f := range d.Functions {
		d.DeploymentMap.Nodes[f.UUID] = f
	}

	fmt.Println("------")
	fmt.Println("Print all edges")

	for i, f := range d.DeploymentMap.Edges {
		fmt.Println("Edge")
		fmt.Println(i)
		fmt.Println(f)
	}
	fmt.Println("------")

	sourceNode, _ := d.DeploymentMap.Nodes[source.UUID]
	var visited map[string]bool
	var stack map[string]bool
	if result := d.IsCyclic(sourceNode, visited, stack); result == true {
		return fmt.Errorf("bad dag, cannot have cycles")
	}
	return nil
}

func (d *DeploymentSpec) AddEdge(FromUUID, ToUUID string, s StreamSpec) error {
	if d.DeploymentMap.Edges[FromUUID] == nil {
		d.DeploymentMap.Edges[FromUUID] = make(map[string]interface{})
	}
	for _, to := range d.DeploymentMap.Edges[FromUUID] {
		if to == ToUUID {
			return fmt.Errorf("Edge exists")
		}
	}
	d.DeploymentMap.Edges[FromUUID][ToUUID] = s

	return nil
}

func (d *DeploymentSpec) IsCyclic(node interface{}, visited, stack map[string]bool) bool {
	fmt.Println("Visiting node")
	fmt.Println(node)
	fmt.Println(visited)
	fmt.Println(stack)

	var nodeID string
	if id, ok := node.(IDInterface); ok {
		fmt.Println(id.ID())
		nodeID = id.ID()
	}
	if visited == nil {
		visited = map[string]bool{}
	}
	if stack == nil {
		stack = map[string]bool{}
	}
	visited[nodeID] = true
	stack[nodeID] = true

	for i, nb := range d.DeploymentMap.Edges[nodeID] {
		fmt.Println("Child node")
		fmt.Println(d.DeploymentMap.Nodes[i])
		search := d.DeploymentMap.Nodes[i]
		if visited[search.(IDInterface).ID()] == false {
			if d.IsCyclic(search, visited, stack) == true {
				return true
			}
		} else if stack[nb.(IDInterface).ID()] == true {
			return true
		}
	}

	stack[nodeID] = false
	return false
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
