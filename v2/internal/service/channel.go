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

func (c *Channel) GracefulStop() {
	c.Stop()
}

func (c *Channel) Stop() {
	c.Channel = nil
}
