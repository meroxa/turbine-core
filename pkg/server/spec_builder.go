package server

import (
	"context"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strings"

	"github.com/google/uuid"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core/v2"
	ir "github.com/meroxa/turbine-core/pkg/ir/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ pb.TurbineServiceServer = (*specBuilderService)(nil)

type specBuilderService struct {
	pb.UnimplementedTurbineServiceServer

	spec *ir.DeploymentSpec
	//resources []*pb.Resource
}

func NewSpecBuilderService() *specBuilderService {
	return &specBuilderService{
		spec: &ir.DeploymentSpec{
			Secrets: make(map[string]string),
		},
	}
}

func (s *specBuilderService) Init(_ context.Context, req *pb.InitRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	s.spec.Definition = ir.DefinitionSpec{
		GitSha: req.GetGitSHA(),
		Metadata: ir.MetadataSpec{
			Turbine: ir.TurbineSpec{
				Language: ir.Lang(strings.ToLower(req.Language.String())),
				Version:  req.TurbineVersion,
			},
			SpecVersion: ir.LatestSpecVersion,
		},
	}
	return empty(), nil
}

func (s *specBuilderService) AddSource(_ context.Context, req *pb.AddSourceRequest) (*pb.AddSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	c := ir.ConnectorSpec{
		UUID:         uuid.New().String(),
		Name:         req.Name,
		PluginType:   ir.PluginSource,
		PluginName:   req.Plugin.Name,
		PluginConfig: configMap(req.Plugin.Configs),
	}

	if err := s.spec.AddSource(&c); err != nil {
		return nil, err
	}

	return &pb.AddSourceResponse{StreamName: c.UUID}, nil
}

func (s *specBuilderService) ReadRecords(_ context.Context, req *pb.ReadRecordsRequest) (*pb.ReadRecordsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &pb.ReadRecordsResponse{
		StreamRecords: &pb.StreamRecords{
			StreamName: req.SourceStream,
		},
	}, nil
}

func (s *specBuilderService) AddDestination(_ context.Context, req *pb.AddDestinationRequest) (*pb.AddDestinationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	c := ir.ConnectorSpec{
		UUID:         uuid.New().String(),
		Name:         req.Name,
		PluginType:   ir.PluginDestination,
		PluginName:   req.Plugin.Name,
		PluginConfig: configMap(req.Plugin.Configs),
	}

	if err := s.spec.AddDestination(&c); err != nil {
		return nil, err
	}

	return &pb.AddDestinationResponse{StreamName: c.UUID}, nil
}

func (s *specBuilderService) WriteRecords(_ context.Context, req *pb.WriteRecordsRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if err := s.spec.AddStream(&ir.StreamSpec{
		UUID:     uuid.New().String(),
		FromUUID: req.StreamRecords.StreamName,
		ToUUID:   req.DestinationID,
		Name:     req.StreamRecords.StreamName + "_" + req.DestinationID,
	}); err != nil {
		return nil, err
	}

	return empty(), nil
}

func (s *specBuilderService) ProcessRecords(_ context.Context, req *pb.ProcessRecordsRequest) (*pb.ProcessRecordsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	f := ir.FunctionSpec{
		UUID: uuid.New().String(),
		Name: strings.ToLower(req.Process.Name),
	}
	if err := s.spec.AddFunction(&f); err != nil {
		return nil, err
	}

	if err := s.spec.AddStream(&ir.StreamSpec{
		UUID:     uuid.New().String(),
		FromUUID: req.StreamRecords.StreamName,
		ToUUID:   f.UUID,
		Name:     req.StreamRecords.StreamName + "_" + f.UUID,
	}); err != nil {
		return nil, err
	}

	return &pb.ProcessRecordsResponse{
		StreamRecords: &pb.StreamRecords{
			StreamName: f.UUID,
			Records:    req.StreamRecords.Records,
		},
	}, nil
}

func (s *specBuilderService) GetSpec(_ context.Context, req *pb.GetSpecRequest) (*pb.GetSpecResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if err := s.spec.SetImageForFunctions(req.Image); err != nil {
		return nil, err
	}

	if _, err := s.spec.BuildDAG(); err != nil {
		return nil, err
	}

	spec, err := s.spec.Marshal()
	if err != nil {
		return nil, err
	}

	return &pb.GetSpecResponse{Spec: spec}, nil
}

func (s *specBuilderService) HasFunctions(_ context.Context, _ *emptypb.Empty) (*wrapperspb.BoolValue, error) {
	return wrapperspb.Bool(len(s.spec.Functions) > 0), nil
}

func configMap(configs *pb.Configs) map[string]any {
	if configs == nil {
		return nil
	}

	m := make(map[string]any)
	for _, c := range configs.Config {
		m[c.Field] = c.Value
	}
	return m
}
