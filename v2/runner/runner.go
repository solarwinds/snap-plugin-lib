package runner

import (
	"github.com/librato/snap-plugin-lib-go/v2/internal/context_manager"
	"github.com/librato/snap-plugin-lib-go/v2/internal/rpc"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

func StartCollector(collector plugin.Collector, name string, version string) {
	contextManager := context_manager.NewContextManager(collector, name, version)
	rpc.StartGRPCController(contextManager)
}
