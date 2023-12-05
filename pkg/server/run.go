package server

import (
	"context"
	"fmt"
	"path"

	pb "github.com/meroxa/turbine-core/v2/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/v2/pkg/app"
	"github.com/meroxa/turbine-core/v2/pkg/server/internal"
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

func (s *runService) AddSource(ctx context.Context, req *pb.AddSourceRequest) (*pb.AddSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &pb.AddSourceResponse{
		StreamName: req.Name,
	}, nil
}

func (s *runService) ReadRecords(ctx context.Context, req *pb.ReadRecordsRequest) (*pb.ReadRecordsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	fixtureFile, ok := s.config.Fixtures[req.SourceStream]
	if !ok {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf(
				"no fixture file found for source %s. Ensure that the source is declared in your app.json.",
				req.SourceStream,
			),
		)
	}

	rr, err := internal.ReadFixture(ctx, path.Join(s.appPath, fixtureFile))
	if err != nil {
		return nil, err
	}

	return &pb.ReadRecordsResponse{
		StreamRecords: &pb.StreamRecords{
			StreamName: req.SourceStream,
			Records:    rr,
		},
	}, nil
}

func (s *runService) AddDestination(ctx context.Context, req *pb.AddDestinationRequest) (*pb.AddDestinationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &pb.AddDestinationResponse{
		StreamName: req.Name,
	}, nil
}

func (s *runService) WriteRecords(ctx context.Context, req *pb.WriteRecordsRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	internal.PrintRecords(req.DestinationID, req.StreamRecords)

	return empty(), nil
}

func (s *runService) ProcessRecords(ctx context.Context, req *pb.ProcessRecordsRequest) (*pb.ProcessRecordsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &pb.ProcessRecordsResponse{
		StreamRecords: &pb.StreamRecords{
			StreamName: req.StreamRecords.StreamName,
			// Records will come from the processing function in the SDK and not the gRPC server
		},
	}, nil
}
