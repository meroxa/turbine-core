package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/servers/info"
	"github.com/meroxa/turbine-core/servers/local"
	"github.com/meroxa/turbine-core/servers/platform"
	"google.golang.org/grpc"
)

const Port = 50051

var (
	Mode string
)

func main() {
	flag.StringVar(&Mode, "mode", "info", "gRPC server mode. Options are info, platform and local. Default is info.")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	log.Printf("running gRPC server in %s mode", Mode)
	switch Mode {
	case "info":
		pb.RegisterTurbineServiceServer(s, info.New())
	case "platform":
		pb.RegisterTurbineServiceServer(s, platform.New())
	case "local":
		pb.RegisterTurbineServiceServer(s, local.New())
	default:
		log.Fatalf("unsupported or invalid mode %s", Mode)
	}
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
