package pluginrpc

import (
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
)

type CollectorProxy interface {
	RequestCollect(id int) ([]*types.Metric, error)
	LoadTask(id int, rawConfig []byte, mtsSelectors []string) error
	UnloadTask(id int) error
	RequestInfo()
}
