package rpc

import (
	"github.com/librato/snap-plugin-lib-go/v2/proxy"
	"golang.org/x/net/context"
)

type collectService struct {
	proxy proxy.Collector
}

func newCollectService(proxy proxy.Collector) CollectorServer {
	return &collectService{
		proxy: proxy,
	}
}

func (cs *collectService) Collect(request *CollectRequest, stream Collector_CollectServer) error {
	log.Trace("GRPC Collect() received")

	taskId := int(request.GetTaskId())

	cs.proxy.RequestCollect(taskId)

	_ = stream.Send(&CollectResponse{
		MetricSet: nil,
	})

	return nil
}

func (cs *collectService) Load(ctx context.Context, request *LoadRequest) (*LoadResponse, error) {
	log.Trace("GRPC Load() received")

	taskId := int(request.GetTaskId())
	jsonConfig := request.GetJsonConfig()
	metrics := request.GetMetricSelector()

	cs.proxy.LoadTask(taskId, jsonConfig, metrics)

	return &LoadResponse{}, nil
}

func (cs *collectService) Unload(ctx context.Context, request *UnloadRequest) (*UnloadResponse, error) {
	log.Trace("GRPC Unload() received")

	taskId := int(request.GetTaskId())

	cs.proxy.UnloadTask(taskId)

	return &UnloadResponse{}, nil
}

func (cs *collectService) Info(ctx context.Context, request *InfoRequest) (*InfoResponse, error) {
	log.Trace("GRPC nfo() received")

	cs.proxy.RequestInfo()

	return &InfoResponse{}, nil
}
