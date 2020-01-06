/*
The package "runner" provides simple API to start plugins in different modes.
*/
package runner

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fullstorydev/grpchan"
	"github.com/sirupsen/logrus"

	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/collector/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/stats"
	"github.com/librato/snap-plugin-lib-go/v2/internal/service"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

var log = logrus.WithFields(logrus.Fields{"layer": "lib", "module": "plugin-runner"})

const (
	normalExitStatus = 0
	errorExitStatus  = 1
)

// As a regular process
func StartCollector(collector plugin.Collector, name string, version string) {
	opt, err := ParseCmdLineOptions(os.Args[0], types.PluginTypeCollector, os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured during plugin startup (%v)\n", err)
		os.Exit(errorExitStatus)
	}

	startCollector(collector, name, version, opt, nil)
}

// As goroutine
func StartCollectorInProcess(publisher plugin.Collector, name string, version string, opt *plugin.Options, grpcChan chan<- grpchan.Channel) {
	startCollector(publisher, name, version, opt, grpcChan)
}

func startCollector(collector plugin.Collector, name string, version string, opt *plugin.Options, grpcChan chan<- grpchan.Channel) {
	var err error

	err = ValidateOptions(opt)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Invalid plugin options (%v)\n", err)
		os.Exit(errorExitStatus)
	}

	statsController, err := stats.NewController(name, version, types.PluginTypeCollector, opt)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error occured when starting statistics controller (%v)\n", err)
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
			_, _ = fmt.Fprintf(os.Stderr, "Can't acquire resources for plugin services (%v)\n", err)
			os.Exit(errorExitStatus)
		}

		printMetaInformation(name, version, types.PluginTypeCollector, opt, r, ctxMan.TasksLimit, ctxMan.InstancesLimit)

		if opt.EnableProfiling {
			startPprofServer(r.pprofListener)
			defer r.pprofListener.Close() // close pprof service when GRPC service has been shut down
		}

		if opt.EnableStatsServer {
			startStatsServer(r.statsListener, statsController)
			defer r.statsListener.Close() // close stats service when GRPC service has been shut down
		}

		srv := service.NewGRPCServer(opt.AsThread)

		// We need to bind the gRPC client on the other end to the same channel so need to return it from here
		if grpcChan != nil {
			grpcChan <- srv.(*service.Channel).Channel
		}

		// main blocking operation
		service.StartCollectorGRPC(srv, ctxMan, statsController, r.grpcListener, r.pprofListener, opt.GRPCPingTimeout, opt.GRPCPingMaxMissed)

	case true:
		startCollectorInSingleMode(ctxMan, opt)
	}
}

func startCollectorInSingleMode(ctxManager *proxy.ContextManager, opt *plugin.Options) {
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
