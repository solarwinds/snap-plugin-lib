package pluginrpc

import "github.com/librato/snap-plugin-lib-go/v2/internal/util/types"

type Publisher interface {
	RequestPublish(id string, mts []*types.Metric) error
	LoadTask(id string, config []byte) error
	UnloadTask(id string) error
}
