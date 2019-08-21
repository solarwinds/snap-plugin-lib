package pluginrpc

import (
	"errors"
	"fmt"

	"github.com/librato/snap-plugin-lib-go/v2/internal/stats"
	"golang.org/x/net/context"
)

const (
	maxCollectChunkSize = 100
)

type collectService struct {
	proxy           CollectorProxy
	statsController stats.Controller
}

func newCollectService(proxy CollectorProxy, statsController stats.Controller) CollectorServer {
	return &collectService{
		proxy:           proxy,
		statsController: statsController,
	}
}

func (cs *collectService) Collect(request *CollectRequest, stream Collector_CollectServer) error {
	log.Trace("GRPC Collect() received")

	taskID := int(request.GetTaskId())

	pluginMts, err := cs.proxy.RequestCollect(taskID)
	if err != nil {
		return fmt.Errorf("plugin is not able to collect metrics: %s", err)
	}

	protoMts := make([]*Metric, 0, len(pluginMts))
	for i, pluginMt := range pluginMts {
		protoMt, err := toGRPCMetric(pluginMt)
		if err != nil {
			log.WithError(err).WithField("metric", pluginMt.Namespace).Errorf("can't send metric over GRPC")
		}

		protoMts = append(protoMts, protoMt)

		if len(protoMts) == maxCollectChunkSize || i == len(pluginMts)-1 {
			err = stream.Send(&CollectResponse{
				MetricSet: protoMts,
			})
			if err != nil {
				log.WithError(err).Error("can't send metric chunk over GRPC")
				return err
			}

			log.WithField("len", len(protoMts)).Debug("metrics chunk has been sent to snap")
		}
	}

	return nil
}

func (cs *collectService) Load(ctx context.Context, request *LoadRequest) (*LoadResponse, error) {
	log.Trace("GRPC Load() received")

	taskID := int(request.GetTaskId())
	jsonConfig := request.GetJsonConfig()
	metrics := request.GetMetricSelectors()

	return &LoadResponse{}, cs.proxy.LoadTask(taskID, jsonConfig, metrics)
}

func (cs *collectService) Unload(ctx context.Context, request *UnloadRequest) (*UnloadResponse, error) {
	log.Trace("GRPC Unload() received")

	taskID := int(request.GetTaskId())

	return &UnloadResponse{}, cs.proxy.UnloadTask(taskID)
}

func (cs *collectService) Info(ctx context.Context, request *InfoRequest) (*InfoResponse, error) {
	log.Trace("GRPC Info() received")

	response := &InfoResponse{}

	select {
	case statistics := <-cs.statsController.RequestStat():
		response.Info = toGRPCInfo(statistics)
	case <-ctx.Done():
		return response, errors.New("couldn't retrieve statistics") // todo: check how this work
	}

	return response, nil
}
