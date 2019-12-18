package runner

import (
	"fmt"
	"net"
)

type resources struct {
	grpcListener  net.Listener
	pprofListener net.Listener
	statsListener net.Listener
}

func (r *resources) grpcListenerAddr() net.TCPAddr {
	if r.grpcListener != nil {
		return *r.grpcListener.Addr().(*net.TCPAddr)
	}

	return net.TCPAddr{}
}

func (r *resources) pprofListenerAddr() net.TCPAddr {
	if r.pprofListener != nil {
		return *r.pprofListener.Addr().(*net.TCPAddr)
	}

	return net.TCPAddr{}
}

func (r *resources) statsListenerAddr() net.TCPAddr {
	if r.statsListener != nil {
		return *r.statsListener.Addr().(*net.TCPAddr)
	}

	return net.TCPAddr{}
}

func acquireResources(opt *Options) (*resources, error) {
	r := &resources{}
	var err error

	r.grpcListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", opt.PluginIP, opt.GRPCPort))
	if err != nil {
		return nil, fmt.Errorf("can't create tcp connection for GRPC server (%s)", err)
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
