package rpc

import (
	"github.com/librato/snap-plugin-lib-go/v2/internal/context_manager"
	"golang.org/x/net/context"
)

type collectService struct {
	proxy context_manager.Collector
}

func newCollectService(proxy context_manager.Collector) CollectorServer {
	return &collectService{
		proxy: proxy,
	}
}

func (cs *collectService) Collect(request *CollectRequest, stream Collector_CollectServer) error {
	log.Trace("GRPC Collect() received")

	taskID := int(request.GetTaskId())

	cs.proxy.RequestCollect(taskID)

	_ = stream.Send(&CollectResponse{
		MetricSet: nil,
	})

	return nil
}

func (cs *collectService) Load(ctx context.Context, request *LoadRequest) (*LoadResponse, error) {
	log.Trace("GRPC Load() received")

	taskID := int(request.GetTaskId())
	jsonConfig := request.GetJsonConfig()
	metrics := request.GetMetricSelector()

	cs.proxy.LoadTask(taskID, jsonConfig, metrics)

	return &LoadResponse{}, nil
}

func (cs *collectService) Unload(ctx context.Context, request *UnloadRequest) (*UnloadResponse, error) {
	log.Trace("GRPC Unload() received")

	taskID := int(request.GetTaskId())

	cs.proxy.UnloadTask(taskID)

	return &UnloadResponse{}, nil
}

func (cs *collectService) Info(ctx context.Context, request *InfoRequest) (*InfoResponse, error) {
	log.Trace("GRPC Info() received")

	cs.proxy.RequestInfo()

	return &InfoResponse{}, nil
}
