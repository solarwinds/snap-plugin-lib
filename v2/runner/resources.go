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
	"fmt"
	"net"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

type resources struct {
	grpcListener  net.Listener
	pprofListener net.Listener
	statsListener net.Listener
}

func safeListenerAddr(ln net.Listener) net.TCPAddr {
	if ln != nil {
		return *ln.Addr().(*net.TCPAddr)
	}

	return net.TCPAddr{}
}

func (r *resources) grpcListenerAddr() net.TCPAddr {
	return safeListenerAddr(r.grpcListener)
}

func (r *resources) pprofListenerAddr() net.TCPAddr {
	return safeListenerAddr(r.pprofListener)
}

func (r *resources) statsListenerAddr() net.TCPAddr {
	return safeListenerAddr(r.statsListener)
}

func acquireResources(opt *plugin.Options) (*resources, error) {
	var err error
	r := &resources{}

	if !opt.AsThread && !opt.DebugMode {
		r.grpcListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", opt.PluginIP, opt.GRPCPort))
		if err != nil {
			return nil, fmt.Errorf("can't create tcp connection for GRPC server (%s)", err)
		}
	}

	if opt.AsThread {
		// force disable profiling as plugin running as goroutine inside Snap will be covered by its pprof anyway
		opt.EnableProfiling = false
	}

	if opt.EnableProfiling {
		r.pprofListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", opt.PluginIP, opt.PProfPort))
		if err != nil {
			return nil, fmt.Errorf("can't create tcp connection for PProf server (%s)", err)
		}
	}

	if opt.EnableStatsServer {
		r.statsListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", opt.PluginIP, opt.StatsPort))
		if err != nil {
			return nil, fmt.Errorf("can't create tcp connection for Stats server (%s)", err)
		}
	}

	return r, nil
}
