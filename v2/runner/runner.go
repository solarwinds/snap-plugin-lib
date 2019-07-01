/*
The package "runner" provides simple API to start plugins in different modes.
*/
package runner

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
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

func StartCollector(collector plugin.Collector, name string, version string) {
	opt, err := ParseCmdLineOptions(os.Args[0], os.Args[1:])
	if err != nil {
		fmt.Printf("Error occured during plugin startup (%v)", err)
		os.Exit(1)
	}

	contextManager := proxy.NewContextManager(collector, name, version)

	logrus.SetLevel(opt.LogLevel)

	switch opt.DebugMode {
	case false:
		r, err := acquireResources(opt)
		if err != nil {
			fmt.Printf("Can't acquire resources for plugin services (%v)", err)
			os.Exit(1)
		}

		printMetaInformation(opt, r)
		startCollectorInServerMode(contextManager, r, opt)
	case true:
		startCollectorInSingleMode(contextManager, opt)
	}
}

func startCollectorInServerMode(ctxManager *proxy.ContextManager, r *resources, opt *plugin.Options) {
	if opt.EnablePprof {
		startPprofServer(opt, r.pprofListener)
	}

	if opt.EnableStats {
		startStatsServer(opt)
	}

	pluginrpc.StartGRPCController(ctxManager, r.grpcListener, opt)
}

func startCollectorInSingleMode(ctxManager *proxy.ContextManager, opt *plugin.Options) {
	const singleModeTaskID = 1

	// Load task based on command line options
	errLoad := ctxManager.LoadTask(singleModeTaskID, []byte(opt.PluginConfig), []string{}) // todo: change this
	if errLoad != nil {
		fmt.Printf("Couldn't load a task in a standalone mode (reason: %v)", errLoad)
		os.Exit(1)
	}

	for runCount := 0; ; {
		// Request metrics collection
		mts, errColl := ctxManager.RequestCollect(singleModeTaskID)
		if errColl != nil {
			fmt.Printf("Error occurred during metrics collection in a standalone mode (reason: %v)", errColl)
			os.Exit(1)
		}

		// Print out metrics
		fmt.Printf("Gathered metrics (length=%d): \n", len(mts))
		for _, mt := range mts {
			fmt.Printf("%s\n", mt) // todo: format output string
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
		os.Exit(1)
	}
}

func printMetaInformation(opt *plugin.Options, r *resources) {
	ip := r.grpcListener.Addr().(*net.TCPAddr).IP.String()

	// Gather meta information
	m := plugin.Meta{
		GRPCVersion: pluginrpc.GRPCDefinitionVersion,
	}

	m.Plugin.Name = ""    // todo: plugin name
	m.Plugin.Version = "" // todo: plugin version

	m.GRPC.IP = ip
	m.GRPC.Port = r.grpcListener.Addr().(*net.TCPAddr).Port

	m.PProf.Enabled = opt.EnablePprof
	if opt.EnablePprof {
		m.PProf.IP = ip
		m.PProf.Port = r.pprofListener.Addr().(*net.TCPAddr).Port
	}

	m.Stats.Enabled = opt.EnableStats
	if opt.EnableStats {
		m.Stats.IP = opt.GrpcIp
		m.Stats.Port = 0 // todo: stats port
	}

	// Print
	jsonMeta, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Can't provide plugin metadata information (reason: %v)\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", string(jsonMeta))
}

func acquireResources(opt *plugin.Options) (*resources, error) {
	r := &resources{}
	var err error

	r.grpcListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", opt.GrpcIp, opt.GrpcPort))
	if err != nil {
		return nil, fmt.Errorf("can't create tcp connection for GRPC server (%s)", err)
	}

	if opt.EnablePprof {
		r.pprofListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", opt.GrpcIp, opt.PprofPort))
		if err != nil {
			return nil, fmt.Errorf("can't create tcp connection for PProf server (%s)", err)
		}
	}

	return r, nil
}
