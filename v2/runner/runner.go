/*
The package "runner" provides simple API to start plugins in different modes.
*/
package runner

import (
	"os"

	"github.com/librato/snap-plugin-lib-go/v2/internal/pluginrpc"
	"github.com/librato/snap-plugin-lib-go/v2/internal/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

func StartCollector(collector plugin.Collector, name string, version string) {
	_, err := ParseCmdLineOptions(os.Args[0], os.Args[1:]) // todo - pass to grpc
	if err != nil {
		os.Exit(1) // todo: more descriptive info
	}

	contextManager := proxy.NewContextManager(collector, name, version)
	pluginrpc.StartGRPCController(contextManager)
}
