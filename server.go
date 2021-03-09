package main

import (
	"context"

	pb "github.com/datachainlab/fabric-ibc-lightclientd/types"
)

type server struct {
	pb.UnimplementedLightClientServer
}

var _ pb.LightClientServer = (*server)(nil)

func (server) CheckHeaderAndUpdateState(ctx context.Context, req *pb.CheckHeaderAndUpdateStateRequest) (*pb.CheckHeaderAndUpdateStateResponse, error) {
	return nil, nil
}
