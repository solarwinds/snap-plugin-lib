package runner

import (
	"github.com/librato/snap-plugin-lib-go/v2/context_manager"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/rpc"
)

func StartCollector(collector plugin.Collector, name string, version string) {

	contextManager := context_manager.NewContextManager(collector, name, version)
	rpc.StartGRPCController(contextManager)
}
