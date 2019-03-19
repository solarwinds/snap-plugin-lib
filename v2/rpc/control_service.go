package rpc

import (
	"errors"
	"golang.org/x/net/context"
)

type controlService struct {
	closeCh chan error
}

func newControlService(closeCh chan error) *controlService {
	return &controlService{
		closeCh: closeCh,
	}
}

func (*controlService) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	log.Trace("GRPC Ping() received")
	return &PingResponse{}, nil
}

func (cs *controlService) Kill(context.Context, *KillRequest) (*KillResponse, error) {
	log.Trace("GRPC Kill() received")

	cs.closeCh <- errors.New("kill requested")
	return &KillResponse{}, nil
}
