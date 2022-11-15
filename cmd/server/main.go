package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/meroxa/turbine-core/pkg/server"
)

const Port = 50051

var (
	Mode string
)

func main() {
	flag.StringVar(&Mode, "mode", "record", "gRPC server mode. Options are record and run. Default is record.")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var s *server.TurbineCoreServer
	log.Printf("running gRPC server in %s mode", Mode)
	switch Mode {
	case "record":
		s = server.NewRecordServer()
	case "run":
		s = server.NewRunServer()
	default:
		log.Fatalf("unsupported or invalid mode %s", Mode)
	}
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
