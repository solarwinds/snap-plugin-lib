/*
 Copyright (c) 2020 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/solarwinds/snap-plugin-lib/v2/internal/service"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/types"
	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
)

// Structure contains information about running services (used by snap)
type meta struct {
	Meta struct {
		RPCVersion string // Message definition version (ie. 2.0.0)
	}

	Plugin struct {
		Name    string           // Plugin name
		Version string           // Plugin version
		Type    types.PluginType // Plugin type (collector, publisher, etc.)
	}

	GRPC struct {
		IP         string // IP on which GRPC service is being served
		Port       int    // Port on which GRPC service is being served
		TLSEnabled bool   // true if TLS is enabled
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

func metaInformation(ctx context.Context, name string, version string, typ types.PluginType, opt *plugin.Options, r *resources, tasksLimit, instancesLimit int) []byte {
	logF := logger(ctx).WithField("service", "meta")

	ip := r.grpcListenerAddr().IP.String()

	m := meta{}

	m.Meta.RPCVersion = service.GRPCDefinitionVersion

	m.Plugin.Name = name
	m.Plugin.Version = version
	m.Plugin.Type = typ

	m.GRPC.IP = ip
	m.GRPC.Port = r.grpcListenerAddr().Port
	m.GRPC.TLSEnabled = opt.EnableTLS

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
		logF.WithError(err).Error("Can't provide plugin metadata information")
		os.Exit(errorExitStatus)
	}

	fmt.Printf("%s\n", string(jsonMeta))
	return jsonMeta
}
