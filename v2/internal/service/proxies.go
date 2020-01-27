package service

import "github.com/librato/snap-plugin-lib-go/v2/internal/util/types"

type CollectorProxy interface {
	RequestCollect(id string) ([]*types.Metric, types.ProcessingError)
	LoadTask(id string, rawConfig []byte, mtsSelectors []string) error
	UnloadTask(id string) error
}
type PublisherProxy interface {
	RequestPublish(id string, mts []*types.Metric) types.ProcessingError
	LoadTask(id string, config []byte) error
	UnloadTask(id string) error
}
