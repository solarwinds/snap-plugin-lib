/*
The package "runner" provides simple API to start plugins in different modes.
*/
package runner

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/librato/snap-plugin-lib-go/v2/internal/pluginrpc"
	"github.com/librato/snap-plugin-lib-go/v2/internal/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"layer": "lib", "module": "plugin-runner"})

func StartCollector(collector plugin.Collector, name string, version string) {
	opt, err := ParseCmdLineOptions(os.Args[0], os.Args[1:])
	if err != nil {
		fmt.Printf("Error occured during plugin startup (%v)", err)
		os.Exit(1)
	}

	standaloneRun := false

	contextManager := proxy.NewContextManager(collector, name, version)
	printMetaInformation(opt)

	logrus.SetLevel(opt.LogLevel)

	switch standaloneRun {
	case false:
		if opt.EnablePprof == true {
			startPprofServer(opt)
		}

		startCollectorInServerMode(contextManager, opt)
	case true:
		startCollectorInSingleMode(contextManager, opt)
	}
}

func startCollectorInServerMode(ctxManager *proxy.ContextManager, opt *plugin.Options) {
	if opt.EnablePprof {
		startPprofServer(opt)
	}

	pluginrpc.StartGRPCController(ctxManager, opt)
}

func startCollectorInSingleMode(ctxManager *proxy.ContextManager, opt *plugin.Options) {
	const singleModeTaskID = 1

	// Load task based on command line options
	errLoad := ctxManager.LoadTask(singleModeTaskID, []byte(opt.PluginConfig), []string{}) // todo: change this
	if errLoad != nil {
		fmt.Printf("Couldn't load a task in a standalone mode (reason: %v)", errLoad)
		os.Exit(1)
	}

	// Request metrics collection
	mts, errColl := ctxManager.RequestCollect(singleModeTaskID)
	if errColl != nil {
		fmt.Printf("Error occurred during metrics collection in a standalone mode (reason: %v)", errColl)
		os.Exit(1)
	}

	// Print out metrics
	fmt.Printf("Gathered metrics (length=%d): \n\n", len(mts))
	for _, mt := range mts {
		fmt.Printf("%#v\n", mt) // todo: format output string
	}

	errUnload := ctxManager.UnloadTask(singleModeTaskID)
	if errUnload != nil {
		fmt.Printf("Couldn't unload a task in a standalone mode (reason: %v)", errUnload)
		os.Exit(1)
	}
}

func printMetaInformation(opt *plugin.Options) {
	// Gather meta information
	m := plugin.Meta{
		GRPCVersion: pluginrpc.GRPCDefinitionVersion,
	}

	m.Plugin.Name = ""    // todo: plugin name
	m.Plugin.Version = "" // todo: plugin version

	m.GRPC.IP = opt.GrpcIp
	m.GRPC.Port = opt.GrpcPort

	m.PProf.Enabled = opt.EnablePprof
	if opt.EnablePprof {
		m.PProf.IP = opt.GrpcIp
		m.PProf.Port = opt.PprofPort
	}

	m.Stats.Enabled = opt.EnableStats
	if opt.EnableStats {
		m.Stats.Enabled = opt.EnableStats
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
