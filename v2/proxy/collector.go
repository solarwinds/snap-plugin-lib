package proxy

import (
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

type Collector interface {
	RequestCollect(id int) ([]plugin.Metric, error)
	LoadTask(id int, config string, selectors []plugin.Namespace) error
	UnloadTask(id int) error
	RequestInfo() plugin.Info
}
