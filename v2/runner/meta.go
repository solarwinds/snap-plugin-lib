package runner

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/librato/snap-plugin-lib-go/v2/internal/pluginrpc"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
)

type PluginType int

const (
	PluginTypeCollector PluginType = iota
	PluginTypeProcessor
	PluginTypePublisher
	PluginTypeStreamingCollector
)

// Structure contains information about running services (used by snap)
type meta struct {
	GRPCVersion string // Message definition version (ie. 2.0.0)

	Plugin struct {
		Name    string     // Plugin name
		Version string     // Plugin version
		Type    PluginType // Plugin type (collector, publisher, etc.)
	}

	GRPC struct {
		IP   string // IP on which GRPC service is being served
		Port int    // Port on which GRPC service is being served
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

func printMetaInformation(name string, version string, opt *types.Options, r *resources) {
	ip := r.grpcListener.Addr().(*net.TCPAddr).IP.String()

	m := meta{
		GRPCVersion: pluginrpc.GRPCDefinitionVersion,
	}

	m.Plugin.Name = name
	m.Plugin.Version = version
	m.Plugin.Type = PluginTypeCollector

	m.GRPC.IP = ip
	m.GRPC.Port = r.grpcListener.Addr().(*net.TCPAddr).Port

	m.Profiling.Enabled = opt.EnableProfiling
	if opt.EnableProfiling {
		m.Profiling.Location = fmt.Sprintf("%s:%d", ip, r.pprofListener.Addr().(*net.TCPAddr).Port)
	}

	m.Stats.Enabled = opt.EnableStatsServer
	if opt.EnableStatsServer {
		m.Stats.IP = ip
		m.Stats.Port = r.statsListener.Addr().(*net.TCPAddr).Port
	}

	// Print
	jsonMeta, err := json.Marshal(m)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Can't provide plugin metadata information (reason: %v)\n", err)
		os.Exit(errorExitStatus)
	}

	fmt.Printf("%s\n", string(jsonMeta))
}
