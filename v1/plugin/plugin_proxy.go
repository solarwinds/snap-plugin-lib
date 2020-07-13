/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

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

package plugin

import (
	"context"
	"sync"
	"time"

	"github.com/librato/snap-plugin-lib-go/v1/plugin/rpc"
	log "github.com/sirupsen/logrus"
)

var (
	// Timeout settings

	// PingTimeoutLimit is the number of successively missed ping health
	// checks which must occur before the plugin is stopped
	PingTimeoutLimit = 3
	// PingTimeoutDuration is the duration during which a ping healthcheck
	// should be received
	PingTimeoutDuration = 6 * time.Second
)

type pluginProxyConstructor func(Plugin) *pluginProxy

type pluginProxy struct {
	plugin              Plugin
	PingTimeoutDuration time.Duration
	halt                chan struct{}

	lastPing   time.Time
	lastPingMu sync.RWMutex
}

// pluginProxyCtor refers to function creating a new plugin proxy instance,
// should never be exposed or used outside of tests
var pluginProxyCtor pluginProxyConstructor = defaultPluginProxyCtor

// newPluginProxy delivers a new plugin proxy instance using a constructor
func newPluginProxy(plugin Plugin) *pluginProxy {
	return pluginProxyCtor(plugin)
}

// defaultPluginProxyCtor delivers new plugin instance using default setup
// (e.g.: default plugin timeout)
func defaultPluginProxyCtor(plugin Plugin) *pluginProxy {
	return &pluginProxy{
		plugin:              plugin,
		PingTimeoutDuration: PingTimeoutDuration,
		halt:                make(chan struct{}),
	}
}

func (p *pluginProxy) Ping(ctx context.Context, arg *rpc.Empty) (*rpc.ErrReply, error) {
	p.lastPingMu.Lock()
	p.lastPing = time.Now()
	p.lastPingMu.Unlock()

	log.WithField("timestamp", p.lastPing).Debug("Heartbeat received")

	return &rpc.ErrReply{}, nil
}

func (p *pluginProxy) Kill(ctx context.Context, arg *rpc.KillArg) (*rpc.ErrReply, error) {
	p.halt <- struct{}{}
	return &rpc.ErrReply{}, nil
}

func (p *pluginProxy) GetConfigPolicy(ctx context.Context, arg *rpc.Empty) (*rpc.GetConfigPolicyReply, error) {
	policy, err := p.plugin.GetConfigPolicy()
	if err != nil {
		return nil, err
	}
	return newGetConfigPolicyReply(policy), nil
}

func (p *pluginProxy) HeartbeatWatch() {
	p.lastPingMu.Lock()
	p.lastPing = time.Now()
	p.lastPingMu.Unlock()

	log.Debug("Heartbeat started")

	count := 0
	for {
		p.lastPingMu.RLock()
		sincePing := time.Since(p.lastPing)
		p.lastPingMu.RUnlock()

		if sincePing >= p.PingTimeoutDuration {
			count++
			log.WithFields(log.Fields{
				"check-duration":     p.PingTimeoutDuration,
				"count":              count,
				"ping-timeout-limit": PingTimeoutLimit,
			}).Warn("Heartbeat timeout")
			if count >= PingTimeoutLimit {
				log.Error("heartbeat timeout expired!")
				defer close(p.halt)
				return
			}
		} else {
			count = 0
		}
		time.Sleep(p.PingTimeoutDuration)
	}

}
