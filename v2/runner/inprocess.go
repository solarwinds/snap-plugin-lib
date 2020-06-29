package runner

import (
	"github.com/librato/grpchan"
	"github.com/sirupsen/logrus"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

type inProcessPlugin interface {
	plugin.InProcessPlugin

	Options() *plugin.Options
	GRPCChannel() chan<- grpchan.Channel
	MetaChannel() chan<- []byte
	Logger() logrus.FieldLogger
}
