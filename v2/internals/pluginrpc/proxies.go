package pluginrpc

import "github.com/librato/snap-plugin-lib-go/v2/internals/util/types"

type CollectorProxy interface {
	RequestCollect(id string) ([]*types.Metric, error)
	LoadTask(id string, rawConfig []byte, mtsSelectors []string) error
	UnloadTask(id string) error
}
type PublisherProxy interface {
	RequestPublish(id string, mts []*types.Metric) error
	LoadTask(id string, config []byte) error
	UnloadTask(id string) error
}
