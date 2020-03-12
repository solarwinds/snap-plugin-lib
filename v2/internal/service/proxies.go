package service

import "github.com/librato/snap-plugin-lib-go/v2/internal/util/types"

type CollectorProxy interface {
	RequestCollect(id string) <-chan types.CollectChunk
	LoadTask(id string, rawConfig []byte, mtsSelectors []string) error
	UnloadTask(id string) error
	CustomInfo(id string) ([]byte, error)
}
type PublisherProxy interface {
	RequestPublish(id string, mts []*types.Metric) types.ProcessingStatus
	LoadTask(id string, config []byte) error
	UnloadTask(id string) error
	CustomInfo(id string) ([]byte, error)
}
