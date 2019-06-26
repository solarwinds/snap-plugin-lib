/*
The package "runner" provides simple API to start plugins in different modes.
*/
package runner

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/librato/snap-plugin-lib-go/v2/internal/pluginrpc"
	"github.com/librato/snap-plugin-lib-go/v2/internal/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

func StartCollector(collector plugin.Collector, name string, version string) {
	opt, err := ParseCmdLineOptions(os.Args[0], os.Args[1:])
	if err != nil {
		fmt.Printf("Error occured during plugin startup (%v)", err)
		os.Exit(1)
	}

	logrus.SetLevel(opt.LogLevel)

	contextManager := proxy.NewContextManager(collector, name, version)
	pluginrpc.StartGRPCController(contextManager, opt)
}
