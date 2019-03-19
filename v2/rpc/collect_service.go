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

func (*collectService) Collect(request *CollectRequest, stream Collector_CollectServer) error {
	log.Trace("GRPC Collect() received")

	_ = stream.Send(&CollectResponse{
		MetricSet: nil,
	})

	return nil
}

func (*collectService) Load(ctx context.Context, request *LoadRequest) (*LoadResponse, error) {
	log.Trace("GRPC Load() received")
	return &LoadResponse{}, nil
}

func (*collectService) Unload(ctx context.Context, request *UnloadRequest) (*UnloadResponse, error) {
	log.Trace("GRPC Unload() received")
	return &UnloadResponse{}, nil
}

func (*collectService) Info(ctx context.Context, request *InfoRequest) (*InfoResponse, error) {
	log.Trace("GRPC Info() received")
	return &InfoResponse{}, nil
}
