package runner

import (
	"github.com/librato/snap-plugin-lib-go/v2/plugin/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/plugin/rpc"
	"github.com/librato/snap-plugin-lib-go/v2/plugin/types"
)

func StartCollector(collector types.Collector, name string, version string) {
	contextManager := proxy.NewContextManager(collector, name, version)
	rpc.StartGRPCController(contextManager)
}
