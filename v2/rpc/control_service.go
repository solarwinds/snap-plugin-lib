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
	return nil, nil
}

func (cs *controlService) Kill(context.Context, *KillRequest) (*KillResponse, error) {
	cs.closeCh <- errors.New("Kill")
	return nil, nil
}
