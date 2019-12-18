package runner

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/librato/snap-plugin-lib-go/v2/internal/pluginrpc"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
)

// Structure contains information about running services (used by snap)
type meta struct {
	GRPCVersion string // Message definition version (ie. 2.0.0)

	Plugin struct {
		Name    string           // Plugin name
		Version string           // Plugin version
		Type    types.PluginType // Plugin type (collector, publisher, etc.)
	}

	GRPC struct {
		IP   string // IP on which GRPC service is being served
		Port int    // Port on which GRPC service is being served
	}

	Constraints struct {
		InstancesLimit int // max number of instances of plugin executable that might be run
		TasksLimit     int // max number of tasks that might be handed per instance
	}

	Profiling struct {
		Enabled  bool   // true, if profiling (pprof server) is enabled (started)
		Location string // location with profiling data (IP and Port for pprof)
	}

	Stats struct {
		Enabled bool   // true, if stats server is enabled (started)
		IP      string // IP on which stats service is being served
		Port    int    // Port on which stats service is being served
	}
}

func printMetaInformation(name string, version string, typ types.PluginType, opt *Options, r *resources, tasksLimit, instancesLimit int) {
	ip := r.grpcListenerAddr().IP.String()

	m := meta{
		GRPCVersion: pluginrpc.GRPCDefinitionVersion,
	}

	m.Plugin.Name = name
	m.Plugin.Version = version
	m.Plugin.Type = typ

	m.GRPC.IP = ip
	m.GRPC.Port = r.grpcListenerAddr().Port

	m.Constraints.TasksLimit = tasksLimit
	m.Constraints.InstancesLimit = instancesLimit

	m.Profiling.Enabled = opt.EnableProfiling
	if opt.EnableProfiling {
		m.Profiling.Location = fmt.Sprintf("%s:%d", ip, r.pprofListenerAddr().Port)
	}

	m.Stats.Enabled = opt.EnableStatsServer
	if opt.EnableStatsServer {
		m.Stats.IP = ip
		m.Stats.Port = r.statsListenerAddr().Port
	}

	// Print
	jsonMeta, err := json.Marshal(m)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Can't provide plugin metadata information (reason: %v)\n", err)
		os.Exit(errorExitStatus)
	}

	fmt.Printf("%s\n", string(jsonMeta))
}
