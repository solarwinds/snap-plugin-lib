package proxy

import (
	"github.com/librato/snap-plugin-lib-go/v2/plugin/types"
)

type Collector interface {
	RequestCollect(id int) []types.Metric
	LoadTask(id int, config string, selectors []types.Namespace) error
	UnloadTask(id int) error
	RequestInfo() types.Info
}
