package collector

import (
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/tutorial/07-proxy/collector/proxy"
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

func (s systemCollector) Collect(plugin.CollectContext) error {
	return nil
}
