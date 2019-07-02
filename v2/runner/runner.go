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
	"github.com/librato/snap-plugin-lib-go/v2/internal/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"layer": "lib", "module": "plugin-runner"})

type resources struct {
	grpcListener  net.Listener
	pprofListener net.Listener
}

const (
	errorStatus = 1
)

func StartCollector(collector plugin.Collector, name string, version string) {
	opt, err := ParseCmdLineOptions(os.Args[0], os.Args[1:])
	if err != nil {
		fmt.Printf("Error occured during plugin startup (%v)", err)
		os.Exit(errorStatus)
	}

	contextManager := proxy.NewContextManager(collector, name, version)

	logrus.SetLevel(opt.LogLevel)

	switch opt.DebugMode {
	case false:
		r, err := acquireResources(opt)
		if err != nil {
			fmt.Printf("Can't acquire resources for plugin services (%v)", err)
			os.Exit(errorStatus)
		}

		printMetaInformation(name, version, opt, r)
		startCollectorInServerMode(contextManager, r, opt)
	case true:
		startCollectorInSingleMode(contextManager, opt)
	}
}

func startCollectorInServerMode(ctxManager *proxy.ContextManager, r *resources, opt *options) {
	if opt.EnablePprof {
		startPprofServer(r.pprofListener)
	}

	if opt.EnableStats {
		startStatsServer()
	}

	pluginrpc.StartGRPCController(ctxManager, r.grpcListener, opt.GrpcPingTimeout, opt.GrpcPingMaxMissed)
}

func startCollectorInSingleMode(ctxManager *proxy.ContextManager, opt *options) {
	const singleModeTaskID = 1

	// Load task based on command line options
	filter := []string{}
	if opt.PluginFilter != defaultFilter {
		filter = strings.Split(opt.PluginFilter, filterSeparator)
	}

	errLoad := ctxManager.LoadTask(singleModeTaskID, []byte(opt.PluginConfig), filter)
	if errLoad != nil {
		fmt.Printf("Couldn't load a task in a standalone mode (reason: %v)", errLoad)
		os.Exit(errorStatus)
	}

	for runCount := uint(0); ; {
		// Request metrics collection
		mts, errColl := ctxManager.RequestCollect(singleModeTaskID)
		if errColl != nil {
			fmt.Printf("Error occurred during metrics collection in a standalone mode (reason: %v)", errColl)
			os.Exit(errorStatus)
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
		fmt.Printf("Couldn't unload a task in a standalone mode (reason: %v)", errUnload)
		os.Exit(errorStatus)
	}
}

func acquireResources(opt *options) (*resources, error) {
	r := &resources{}
	var err error

	r.grpcListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", opt.PluginIp, opt.GrpcPort))
	if err != nil {
		return nil, fmt.Errorf("can't create tcp connection for GRPC server (%s)", err)
	}

	if opt.EnablePprof {
		r.pprofListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", opt.PluginIp, opt.PprofPort))
		if err != nil {
			return nil, fmt.Errorf("can't create tcp connection for PProf server (%s)", err)
		}
	}

	return r, nil
}
