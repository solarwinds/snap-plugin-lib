package pluginrpc

import (
	"errors"
	"fmt"
	"net"

	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/collector/stats"
	"golang.org/x/net/context"
)

const (
	maxCollectChunkSize = 100
)

var logCollectService = log.WithField("service", "Collect")

type collectService struct {
	proxy           CollectorProxy
	statsController stats.Controller
	pprofLn         net.Listener
}

func newCollectService(proxy CollectorProxy, statsController stats.Controller, pprofLn net.Listener) CollectorServer {
	return &collectService{
		proxy:           proxy,
		statsController: statsController,
		pprofLn:         pprofLn,
	}
}

func (cs *collectService) Collect(request *CollectRequest, stream Collector_CollectServer) error {
	logCollectService.Debug("GRPC Collect() received")

	taskID := string(request.GetTaskId())

	pluginMts, err := cs.proxy.RequestCollect(taskID)
	if err != nil {
		return fmt.Errorf("plugin is not able to collect metrics: %s", err)
	}

	protoMts := make([]*Metric, 0, len(pluginMts))
	for i, pluginMt := range pluginMts {
		protoMt, err := toGRPCMetric(pluginMt)
		if err != nil {
			logCollectService.WithError(err).WithField("metric", pluginMt.Namespace).Errorf("can't send metric over GRPC")
		}

		protoMts = append(protoMts, protoMt)

		if len(protoMts) == maxCollectChunkSize || i == len(pluginMts)-1 {
			err = stream.Send(&CollectResponse{
				MetricSet: protoMts,
			})
			if err != nil {
				logCollectService.WithError(err).Error("can't send metric chunk over GRPC")
				return err
			}

			logCollectService.WithField("len", len(protoMts)).Debug("metrics chunk has been sent to snap")
		}
	}

	return nil
}

func (cs *collectService) Load(ctx context.Context, request *LoadCollectorRequest) (*LoadCollectorResponse, error) {
	logCollectService.Debug("GRPC Load() received")

	taskID := string(request.GetTaskId())
	jsonConfig := request.GetJsonConfig()
	metrics := request.GetMetricSelectors()

	return &LoadCollectorResponse{}, cs.proxy.LoadTask(taskID, jsonConfig, metrics)
}

func (cs *collectService) Unload(ctx context.Context, request *UnloadCollectorRequest) (*UnloadCollectorResponse, error) {
	logCollectService.Debug("GRPC Unload() received")

	taskID := string(request.GetTaskId())

	return &UnloadCollectorResponse{}, cs.proxy.UnloadTask(taskID)
}

func (cs *collectService) Info(ctx context.Context, request *InfoRequest) (*InfoResponse, error) {
	logCollectService.Debug("GRPC Info() received")

	var err error
	response := &InfoResponse{}

	select {
	case statistics := <-cs.statsController.RequestStat():
		if statistics == nil {
			return response, fmt.Errorf("can't gather statistics (statistics server is not running): %v", err)
		}

		pprofAddr := ""
		if cs.pprofLn != nil {
			pprofAddr = cs.pprofLn.Addr().String()
		}

		response.Info, err = toGRPCInfo(statistics, pprofAddr)
		if err != nil {
			return response, fmt.Errorf("could't convert statistics to GRPC format: %v", err)
		}
	case <-ctx.Done():
		return response, errors.New("won't retrieve statistics - request canceled")
	}

	return response, nil
}
