package runner

import (
	"fmt"
	"github.com/librato/snap-plugin-lib-go/v2/internal/pluginrpc"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"os"
)

func StartPublisher(publisher plugin.Publisher, name string, version string) {
	opt, err := ParseCmdLineOptions(os.Args[0], os.Args[1:]) // todo: publish has limited set of flags
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error occured during plugin startup (%v)\n", err)
		os.Exit(errorExitStatus)
	}

	//contextManager := proxy.NewContextManager(publisher)

	r, err := acquireResources(opt) // todo: but no pprof and stats
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Can't acquire resources for plugin services (%v)\n", err)
		os.Exit(errorExitStatus)
	}

	printMetaInformation(name, version, PluginTypePublisher, opt, r)
	pluginrpc.StartPublisherGRPC(r.grpcListener, opt.GrpcPingTimeout, opt.GrpcPingMaxMissed)
}
