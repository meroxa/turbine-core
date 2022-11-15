package server

import (
	"context"
	"log"
	"net"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	ListenAddress = "localhost:50500"
)

type TurbineCoreServer struct {
	*grpc.Server
}

func NewRunServer() *TurbineCoreServer {
	s := grpc.NewServer()
	pb.RegisterTurbineServiceServer(s, NewRunService())
	return &TurbineCoreServer{Server: s}
}

func NewRecordServer() *TurbineCoreServer {
	s := grpc.NewServer()
	pb.RegisterTurbineServiceServer(s, NewRunService())
	return &TurbineCoreServer{Server: s}
}

func (s *TurbineCoreServer) Run(ctx context.Context) {
	listener, err := net.Listen("tcp", ListenAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func Empty() *emptypb.Empty {
	return new(emptypb.Empty)
}