package service

import (
	"net"

	"github.com/fullstorydev/grpchan/inprocgrpc"
)

type Channel struct {
	*inprocgrpc.Channel
}

func NewChannel() *Channel {
	return &Channel{
		&inprocgrpc.Channel{},
	}
}

func (*Channel) Serve(_ net.Listener) error {
	return nil
}

func (*Channel) GracefulStop() {
	panic("Implement me")
}

func (*Channel) Stop() {
	panic("Implement me")
}
