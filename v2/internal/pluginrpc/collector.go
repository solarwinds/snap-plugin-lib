package pluginrpc

import "github.com/librato/snap-plugin-lib-go/v2/plugin"

type CollectorProxy interface {
	RequestCollect(id int) ([]plugin.Metric, error)
	LoadTask(id int, rawConfig []byte, mtsSelectors []string) error
	UnloadTask(id int) error
	RequestInfo()
}
