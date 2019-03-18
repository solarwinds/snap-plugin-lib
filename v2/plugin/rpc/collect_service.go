package rpc

import (
	"github.com/librato/snap-plugin-lib-go/v2/plugin/proxy"
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

func (*collectService) Collect(*CollectRequest, Collector_CollectServer) error {
	return nil
}

func (*collectService) Load(context.Context, *LoadRequest) (*LoadResponse, error) {
	return nil, nil
}

func (*collectService) Unload(context.Context, *UnloadRequest) (*UnloadResponse, error) {
	return nil, nil
}

func (*collectService) Info(context.Context, *InfoRequest) (*InfoResponse, error) {
	return nil, nil
}
