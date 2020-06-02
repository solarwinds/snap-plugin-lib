package runner

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/stats"
	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/publisher/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/service"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

func StartPublisher(publisher plugin.Publisher, name string, version string) {
	opt, err := ParseCmdLineOptions(os.Args[0], types.PluginTypePublisher, os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error occured during plugin startup (%v)\n", err)
		os.Exit(errorExitStatus)
	}

	startPublisher(publisher, name, version, opt, InProcChanWriteOnly{})
}

func startPublisher(publisher plugin.Publisher, name string, version string, opt *plugin.Options, inProcChan InProcChanWriteOnly) {
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

	if opt.PrintVersion {
		printVersion(name, version)
		os.Exit(normalExitStatus)
	}

	r, err := acquireResources(opt)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Can't acquire resources for plugin services (%v)\n", err)
		os.Exit(errorExitStatus)
	}

	jsonMeta := metaInformation(name, version, types.PluginTypePublisher, opt, r, ctxMan.TasksLimit, ctxMan.InstancesLimit)
	if inProcChan.MetaCh != nil {
		inProcChan.MetaCh <- jsonMeta
		close(inProcChan.MetaCh)
	}

	if opt.EnableProfiling {
		startPprofServer(r.pprofListener)
		defer r.pprofListener.Close() // close pprof service when GRPC service has been shut down
	}

	if opt.EnableStatsServer {
		startStatsServer(r.statsListener, statsController)
		defer r.statsListener.Close() // close stats service when GRPC service has been shut down
	}

	srv, err := service.NewGRPCServer(opt)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Can't initialize GRPC Server (%v)\n", err)
		os.Exit(errorExitStatus)
	}

	// We need to bind the gRPC client on the other end to the same channel so need to return it from here
	if inProcChan.GRPCChan != nil {
		inProcChan.GRPCChan <- srv.(*service.Channel).Channel
	}

	// main blocking operation
	service.StartPublisherGRPC(srv, ctxMan, r.grpcListener, opt.GRPCPingTimeout, opt.GRPCPingMaxMissed)
}
