package runner

import (
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/rpc"
	"github.com/sirupsen/logrus"
)

func StartCollector(collector plugin.Collector, name string, version string) {
	logrus.SetLevel(logrus.TraceLevel) // todo: remove

	contextManager := proxy.NewContextManager(collector, name, version)
	rpc.StartGRPCController(contextManager)
}
