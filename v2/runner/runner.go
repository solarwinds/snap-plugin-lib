package runner

import (
	"github.com/librato/snap-plugin-lib-go/v2/internal/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/rpc"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

func StartCollector(collector plugin.Collector, name string, version string) {
	contextManager := proxy.NewContextManager(collector, name, version)
	rpc.StartGRPCController(contextManager)
}
