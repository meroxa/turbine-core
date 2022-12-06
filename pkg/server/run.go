package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/meroxa/turbine-core/pkg/app"
	"github.com/meroxa/turbine-core/pkg/ir"
	"github.com/meroxa/turbine-core/pkg/server/internal"

	"github.com/google/uuid"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ pb.TurbineServiceServer = (*runService)(nil)

type runService struct {
	pb.UnimplementedTurbineServiceServer

	deploymentSpec ir.DeploymentSpec
	resources      []*pb.Resource

	config  app.Config
	appPath string
}

func NewRunService() *runService {
	return &runService{
		config: app.Config{},
	}
}

func (s *runService) Init(ctx context.Context, req *pb.InitRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	s.deploymentSpec.Definition = ir.DefinitionSpec{}

	config, err := app.ReadConfig(req.AppName, req.ConfigFilePath)
	if err != nil {
		return nil, err
	}
	s.config = config
	s.appPath = req.ConfigFilePath

	return empty(), nil
}

func (s *runService) GetResource(ctx context.Context, req *pb.GetResourceRequest) (*pb.Resource, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &pb.Resource{
		Name: req.Name,
	}, nil
}

func (s *runService) ReadCollection(ctx context.Context, req *pb.ReadCollectionRequest) (*pb.Collection, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	fixture := &internal.FixtureResource{
		Collection: req.Collection,
		File: path.Join(
			s.appPath,
			s.config.Resources[req.Resource.Name],
		),
	}

	rr, err := fixture.ReadAll(ctx)
	if err != nil {
		return nil, err
	}

	// spec
	s.resources = append(s.resources, &pb.Resource{
		Name:       req.GetResource().GetName(),
		Source:     true,
		Collection: req.GetCollection(),
	})

	for _, c := range s.deploymentSpec.Connectors {
		// Only one source per app allowed.
		if c.Type == ir.ConnectorSource {
			return &pb.Collection{}, fmt.Errorf("only one call to 'read' is allowed per Meroxa Data Application")
		}
	}
	c := ir.ConnectorSpec{
		ID:         uuid.New().String(),
		Collection: req.GetCollection(),
		Resource:   req.Resource.GetName(),
		Type:       ir.ConnectorSource,
		Config:     resourceConfigsToMap(req.GetConfigs().GetConfig()),
	}

	s.deploymentSpec.Connectors = append(s.deploymentSpec.Connectors, c)

	return &pb.Collection{
		Name:    req.Collection,
		Records: rr,
		Stream:  c.ID,
	}, nil
}

func (s *runService) WriteCollectionToResource(ctx context.Context, req *pb.WriteCollectionRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	source := req.GetSourceCollection()
	c := ir.ConnectorSpec{
		ID:         uuid.New().String(),
		Collection: req.GetTargetCollection(),
		Resource:   req.Resource.GetName(),
		Type:       ir.ConnectorDestination,
		Config:     resourceConfigsToMap(req.GetConfigs().GetConfig()),
	}
	s.deploymentSpec.Connectors = append(s.deploymentSpec.Connectors, c)

	stream := ir.StreamSpec{
		ID:     uuid.New().String(),
		FromID: source.Stream, // ID of the source node
		ToID:   c.ID,
	}
	s.deploymentSpec.Streams = append(s.deploymentSpec.Streams, stream)

	internal.PrintRecords(
		req.Resource.Name,
		req.TargetCollection,
		req.SourceCollection.Records,
	)

	spec, err := s.deploymentSpec.Marshal()
	if err != nil {
		return empty(), err
	}

	var b bytes.Buffer
	if err := json.Indent(&b, spec, "", "\t"); err != nil {
		return empty(), err
	}
	b.WriteTo(os.Stdout)

	return empty(), nil
}

func (s *runService) AddProcessToCollection(ctx context.Context, req *pb.ProcessCollectionRequest) (*pb.Collection, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	collection := req.GetCollection()
	p := req.GetProcess()

	f := ir.FunctionSpec{
		ID:   uuid.New().String(),
		Name: strings.ToLower(p.GetName()),
	}
	s.deploymentSpec.Functions = append(s.deploymentSpec.Functions, f)

	stream := ir.StreamSpec{
		ID:     uuid.New().String(),
		FromID: collection.Stream, // ID of the source node
		ToID:   f.ID,
	}
	s.deploymentSpec.Streams = append(s.deploymentSpec.Streams, stream)

	return &pb.Collection{
		Name:    "function-records",
		Records: collection.Records,
		Stream:  f.ID,
	}, nil
}

func (s *runService) RegisterSecret(ctx context.Context, req *pb.Secret) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return empty(), nil
}
