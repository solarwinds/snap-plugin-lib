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

	if !opt.AsThread {
		r.grpcListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", opt.PluginIP, opt.GRPCPort))
		if err != nil {
			return nil, fmt.Errorf("can't create tcp connection for GRPC server (%s)", err)
		}
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
