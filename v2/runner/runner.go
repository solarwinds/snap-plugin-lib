/*
The package "runner" provides simple API to start plugins in different modes.
*/
package runner

import (
	"github.com/librato/snap-plugin-lib-go/v2/internal/pluginrpc"
	"github.com/librato/snap-plugin-lib-go/v2/internal/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

func StartCollector(collector plugin.Collector, name string, version string) {
	contextManager := proxy.NewContextManager(collector, name, version)
	pluginrpc.StartGRPCController(contextManager)
}
