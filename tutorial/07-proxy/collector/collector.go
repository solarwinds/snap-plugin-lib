package collector

import (
	"github.com/librato/snap-plugin-lib-go/tutorial/07-proxy/collector/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

type systemCollector struct {
	proxyCollector proxy.Proxy
}

func (s systemCollector) PluginDefinition(def plugin.CollectorDefinition) error {
	return nil
}

func New(proxy proxy.Proxy) systemCollector {
	return systemCollector{
		proxyCollector: proxy,
	}
}

func (s systemCollector) Collect(plugin.Context) error {
	return nil
}
