package server

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/meroxa/turbine-core/pkg/app"
	"github.com/meroxa/turbine-core/pkg/server/internal"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
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
	config, err := app.ReadConfig(req.AppName, req.ConfigFilePath)
	if err != nil {
		return empty(), err
	}
	s.config = config
	s.appPath = req.ConfigFilePath

	return empty(), nil
}

func (s *runService) GetResource(ctx context.Context, id *pb.GetResourceRequest) (*pb.Resource, error) {
	return &pb.Resource{
		Name: id.Name,
	}, nil
}

func (s *runService) ReadCollection(ctx context.Context, request *pb.ReadCollectionRequest) (*pb.Collection, error) {
	if request.Collection == "" {
		return &pb.Collection{}, fmt.Errorf("please provide a collection name to Records()")
	}

	fixtureFile := s.config.Resources[request.Resource.Name]
	resourceFixturesPath := fmt.Sprintf("%s/%s", s.appPath, fixtureFile)
	return internal.ReadFixtures(resourceFixturesPath, request.Collection)
}

func (s *runService) WriteCollectionToResource(ctx context.Context, request *pb.WriteCollectionRequest) (*emptypb.Empty, error) {
	if request.SourceCollection.Name == "" {
		return empty(), fmt.Errorf("please provide a collection name to Records()")
	}

	internal.PrettyPrintRecords(request.Resource.Name, request.SourceCollection.Stream, request.SourceCollection.Records)

	return empty(), nil
}

func (s *runService) AddProcessToCollection(ctx context.Context, request *pb.ProcessCollectionRequest) (*pb.Collection, error) {
	return request.GetCollection(), nil
}

func (s *runService) RegisterSecret(ctx context.Context, secret *pb.Secret) (*emptypb.Empty, error) {
	val := os.Getenv(secret.Name)
	if val == "" {
		return empty(), errors.New("secret is invalid or not set")
	}
	return empty(), nil
}
