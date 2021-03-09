package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/datachainlab/fabric-ibc-lightclientd/types"
)

const address = ":60000"

func main() {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterLightClientServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
