package runner

import (
	"fmt"
	"os"

	"github.com/librato/snap-plugin-lib-go/v2/internal/pluginrpc"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

func StartPublisher(publisher plugin.Publisher, name string, version string) {
	opt, err := ParseCmdLineOptions(os.Args[0], PluginTypePublisher, os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error occured during plugin startup (%v)\n", err)
		os.Exit(errorExitStatus)
	}

	//contextManager := proxy.NewContextManager(publisher)

	r, err := acquireResources(opt)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Can't acquire resources for plugin services (%v)\n", err)
		os.Exit(errorExitStatus)
	}

	printMetaInformation(name, version, PluginTypePublisher, opt, r)
	pluginrpc.StartPublisherGRPC(r.grpcListener, opt.GrpcPingTimeout, opt.GrpcPingMaxMissed)
}
