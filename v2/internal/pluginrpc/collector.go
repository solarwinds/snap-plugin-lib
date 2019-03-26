package pluginrpc

import "github.com/librato/snap-plugin-lib-go/v2/plugin"

type CollectorProxy interface {
	RequestCollect(id int) ([]plugin.Metric, error)
	LoadTask(id int, config string, mtsSelectors []string) error
	UnloadTask(id int) error
	RequestInfo()
}
