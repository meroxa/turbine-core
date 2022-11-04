package server

import (
	"context"
	"log"
	"net"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"

	"google.golang.org/grpc"
)

const (
	ListenAddress = "localhost:50500"
)

type turbineCoreServer struct {
	*grpc.Server
}

func NewRunServer() *turbineCoreServer {
	s := grpc.NewServer()
	pb.RegisterTurbineServiceServer(s, NewRunService())
	return &turbineCoreServer{Server: s}
}

func (s *turbineCoreServer) Run(ctx context.Context) {
	listener, err := net.Listen("tcp", ListenAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
