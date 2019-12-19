package runner

import (
	"fmt"
	"os"

	"github.com/fullstorydev/grpchan"
	"github.com/sirupsen/logrus"

	"github.com/librato/snap-plugin-lib-go/v2/internal/pluginrpc"
	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/stats"
	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/publisher/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

// As a regular process
func StartPublisher(publisher plugin.Publisher, name string, version string) {
	opt, err := ParseCmdLineOptions(os.Args[0], types.PluginTypePublisher, os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error occured during plugin startup (%v)\n", err)
		os.Exit(errorExitStatus)
	}

	startPublisher(publisher, name, version, opt, nil)
}

// As goroutine
func StartPublisherInProcess(publisher plugin.Publisher, name string, version string, opt *plugin.Options, grpcChan chan<- grpchan.Channel) {
	startPublisher(publisher, name, version, opt, grpcChan)
}

func startPublisher(publisher plugin.Publisher, name string, version string, opt *plugin.Options, grpcChan chan<- grpchan.Channel) {
	var err error

	err = ValidateOptions(opt)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Invalid plugin options (%v)\n", err)
		os.Exit(errorExitStatus)
	}

	statsController, err := stats.NewController(name, version, types.PluginTypePublisher, opt)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error occured when starting statistics controller (%v)\n", err)
		os.Exit(errorExitStatus)
	}

	ctxMan := proxy.NewContextManager(publisher, statsController)

	logrus.SetLevel(opt.LogLevel)

	r := &resources{}
	if !opt.AsThread {
		r, err = acquireResources(opt)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Can't acquire resources for plugin services (%v)\n", err)
			os.Exit(errorExitStatus)
		}
	}

	printMetaInformation(name, version, types.PluginTypePublisher, opt, r, ctxMan.TasksLimit, ctxMan.InstancesLimit)

	if opt.EnableProfiling {
		startPprofServer(r.pprofListener)
		defer r.pprofListener.Close() // close pprof service when GRPC service has been shut down
	}

	if opt.EnableStatsServer {
		startStatsServer(r.statsListener, statsController)
		defer r.statsListener.Close() // close stats service when GRPC service has been shut down
	}

	srv := pluginrpc.NewGRPCServer(opt.AsThread)
	if grpcChan != nil {
		grpcChan <- srv.(*pluginrpc.Channel)
	}

	// main blocking operation
	pluginrpc.StartPublisherGRPC(srv, ctxMan, statsController, r.grpcListener, r.pprofListener, opt.GRPCPingTimeout, opt.GRPCPingMaxMissed)
}
