package ir

import (
	"encoding/json"
	"fmt"
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
	Secrets    map[string]string `json:"secrets,omitempty"`
	Connectors []ConnectorSpec   `json:"connectors"`
	Functions  []FunctionSpec    `json:"functions,omitempty"`
	Streams    []StreamSpec      `json:"streams,omitempty"`
	Definition DefinitionSpec    `json:"definition"`
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
