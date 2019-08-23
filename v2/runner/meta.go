package runner

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/librato/snap-plugin-lib-go/v2/internal/pluginrpc"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
)

// Structure contains information about running services (used by snap)
type meta struct {
	GRPCVersion string // Message definition version (ie. 2.0.0)

	Plugin struct {
		Name    string // Plugin name
		Version string // Plugin version
	}

	GRPC struct {
		IP   string // IP on which GRPC service is being served
		Port int    // Port on which GRPC service is being served
	}

	PProf struct {
		Enabled bool   // true, if pprof server is enabled (started)
		IP      string // IP on which pprof service is being served
		Port    int    // Port on which pprof service is being served
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

	m.GRPC.IP = ip
	m.GRPC.Port = r.grpcListener.Addr().(*net.TCPAddr).Port

	m.PProf.Enabled = opt.EnablePprofServer
	if opt.EnablePprofServer {
		m.PProf.IP = ip
		m.PProf.Port = r.pprofListener.Addr().(*net.TCPAddr).Port
	}

	m.Stats.Enabled = opt.EnableStatsServer
	if opt.EnableStatsServer {
		m.Stats.IP = ip
		m.Stats.Port = r.statsListener.Addr().(*net.TCPAddr).Port
	}

	// Print
	jsonMeta, err := json.Marshal(m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't provide plugin metadata information (reason: %v)\n", err)
		os.Exit(errorExitStatus)
	}

	fmt.Printf("%s\n", string(jsonMeta))
}
