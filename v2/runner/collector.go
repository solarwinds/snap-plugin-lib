/*
The package "runner" provides simple API to start plugins in different modes.
*/
package runner

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/pluginrpc"
	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/collector/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/stats"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"layer": "lib", "module": "plugin-runner"})

type resources struct {
	grpcListener  net.Listener
	pprofListener net.Listener
	statsListener net.Listener
}

const (
	normalExitStatus = 0
	errorExitStatus  = 1
)

func StartCollector(collector plugin.Collector, name string, version string) {
	opt, err := ParseCmdLineOptions(os.Args[0], types.PluginTypeCollector, os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured during plugin startup (%v)\n", err)
		os.Exit(errorExitStatus)
	}

	statsController, err := stats.NewController(name, version, types.PluginTypeCollector, opt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured when starting statistics controller (%v)\n", err)
		os.Exit(errorExitStatus)
	}

	ctxMan := proxy.NewContextManager(collector, statsController)

	logrus.SetLevel(opt.LogLevel)

	if opt.PrintExampleTask {
		printExampleTask(ctxMan, name)
		os.Exit(normalExitStatus)
	}

	switch opt.DebugMode {
	case false:
		r, err := acquireResources(opt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't acquire resources for plugin services (%v)\n", err)
			os.Exit(errorExitStatus)
		}

		printMetaInformation(name, version, types.PluginTypeCollector, opt, r, ctxMan.TasksLimit, ctxMan.InstancesLimit)
		startCollectorInServerMode(ctxMan, statsController, r, opt)
	case true:
		startCollectorInSingleMode(ctxMan, opt)
	}
}

func startCollectorInServerMode(ctxManager *proxy.ContextManager, statsController stats.Controller, r *resources, opt *types.Options) {
	if opt.EnableProfiling {
		startPprofServer(r.pprofListener)
		defer r.pprofListener.Close() // close pprof service when GRPC service has been shut down
	}

	if opt.EnableStatsServer {
		startStatsServer(r.statsListener, statsController)
		defer r.statsListener.Close() // close stats service when GRPC service has been shut down
	}

	// main blocking operation
	pluginrpc.StartCollectorGRPC(ctxManager, statsController, r.grpcListener, r.pprofListener, opt.GrpcPingTimeout, opt.GrpcPingMaxMissed)
}

func startCollectorInSingleMode(ctxManager *proxy.ContextManager, opt *types.Options) {
	const singleModeTaskID = "task-1"

	// Load task based on command line options
	filter := []string{}
	if opt.PluginFilter != defaultFilter {
		filter = strings.Split(opt.PluginFilter, filterSeparator)
	}

	errLoad := ctxManager.LoadTask(singleModeTaskID, []byte(opt.PluginConfig), filter)
	if errLoad != nil {
		fmt.Fprintf(os.Stderr, "Couldn't load a task in a standalone mode (reason: %v)\n", errLoad)
		os.Exit(errorExitStatus)
	}

	for runCount := uint(0); ; {
		// Request metrics collection
		mts, errColl := ctxManager.RequestCollect(singleModeTaskID)
		if errColl != nil {
			fmt.Fprintf(os.Stderr, "Error occurred during metrics collection in a standalone mode (reason: %v)\n", errColl)
			os.Exit(errorExitStatus)
		}

		// Print out metrics
		fmt.Printf("Gathered metrics (length=%d): \n", len(mts))
		for _, mt := range mts {
			fmt.Printf("%s\n", mt)
		}
		fmt.Printf("\n")

		// wait to request new collection or exit
		runCount++
		if runCount == opt.DebugCollectCounts {
			break
		}

		time.Sleep(opt.DebugCollectInterval)
	}

	errUnload := ctxManager.UnloadTask(singleModeTaskID)
	if errUnload != nil {
		fmt.Fprintf(os.Stderr, "Couldn't unload a task in a standalone mode (reason: %v)\n", errUnload)
		os.Exit(errorExitStatus)
	}
}

func acquireResources(opt *types.Options) (*resources, error) {
	r := &resources{}
	var err error

	r.grpcListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", opt.PluginIp, opt.GrpcPort))
	if err != nil {
		return nil, fmt.Errorf("can't create tcp connection for GRPC server (%s)", err)
	}

	if opt.EnableProfiling {
		r.pprofListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", opt.PluginIp, opt.PprofPort))
		if err != nil {
			return nil, fmt.Errorf("can't create tcp connection for PProf server (%s)", err)
		}
	}

	if opt.EnableStatsServer {
		r.statsListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", opt.PluginIp, opt.StatsPort))
		if err != nil {
			return nil, fmt.Errorf("can't create tcp connection for Stats server (%s)", err)
		}
	}

	return r, nil
}
