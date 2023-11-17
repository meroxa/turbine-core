package server

import (
	"context"
	"fmt"
	"path"

	"github.com/meroxa/turbine-core/pkg/app"
	"github.com/meroxa/turbine-core/pkg/server/internal"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ pb.TurbineServiceServer = (*runService)(nil)

type runService struct {
	pb.UnimplementedTurbineServiceServer

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

	config, err := app.ReadConfig(req.AppName, req.ConfigFilePath)
	if err != nil {
		return nil, err
	}
	s.config = config
	s.appPath = req.ConfigFilePath

	return empty(), nil
}

func (s *runService) GetSource(ctx context.Context, req *pb.GetSourceRequest) (*pb.Source, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &pb.Source{
		Name: req.Name,
	}, nil
}

func (s *runService) GetDestination(ctx context.Context, req *pb.GetDestinationRequest) (*pb.Destination, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &pb.Destination{
		Name: req.Name,
	}, nil
}

func (s *runService) ReadCollection(ctx context.Context, req *pb.ReadCollectionRequest) (*pb.Collection, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	fixtureFile, ok := s.config.Fixtures[req.Source.Name]
	if !ok {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf(
				"No fixture file found for source %s. Ensure that the source is declared in your app.json.",
				req.Source.Name,
			),
		)
	}

	fixture := &internal.FixtureResource{
		Collection: req.Collection,
		File: path.Join(
			s.appPath,
			fixtureFile,
		),
	}

	rr, err := fixture.ReadAll(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.Collection{
		Name:    req.Collection,
		Records: rr,
	}, nil
}

func (s *runService) WriteCollectionToResource(ctx context.Context, req *pb.WriteCollectionRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	internal.PrintRecords(
		req.Destination.Name,
		req.DestinationCollection,
		req.SourceCollection.Records,
	)

	return empty(), nil
}

func (s *runService) AddProcessToCollection(ctx context.Context, req *pb.ProcessCollectionRequest) (*pb.Collection, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return req.Collection, nil
}

func (s *runService) RegisterSecret(ctx context.Context, req *pb.Secret) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return empty(), nil
}
