package service

import (
	"context"
	"fmt"
	"net"

	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/stats"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/pluginrpc"
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

func newCollectService(proxy CollectorProxy, statsController stats.Controller, pprofLn net.Listener) pluginrpc.CollectorServer {
	return &collectService{
		proxy:           proxy,
		statsController: statsController,
		pprofLn:         pprofLn,
	}
}

func (cs *collectService) Collect(request *pluginrpc.CollectRequest, stream pluginrpc.Collector_CollectServer) error {
	logCollectService.Debug("GRPC Collect() received")

	taskID := request.GetTaskId()

	pluginMts, status := cs.proxy.RequestCollect(taskID)

	err := cs.collectWarnings(stream, status.Warnings)
	if err != nil {
		return fmt.Errorf("can't send all warnings to snap: %v", err)
	}

	if status.Error != nil {
		return fmt.Errorf("plugin is not able to collect metrics: %s", status)
	}

	err = cs.collectMetrics(stream, pluginMts)
	if err != nil {
		return fmt.Errorf("can't send all metrics to snap: %v", err)
	}

	return nil
}

func (cs *collectService) Load(ctx context.Context, request *pluginrpc.LoadCollectorRequest) (*pluginrpc.LoadCollectorResponse, error) {
	logCollectService.Debug("GRPC Load() received")

	taskID := string(request.GetTaskId())
	jsonConfig := request.GetJsonConfig()
	metrics := request.GetMetricSelectors()

	return &pluginrpc.LoadCollectorResponse{}, cs.proxy.LoadTask(taskID, jsonConfig, metrics)
}

func (cs *collectService) Unload(ctx context.Context, request *pluginrpc.UnloadCollectorRequest) (*pluginrpc.UnloadCollectorResponse, error) {
	logCollectService.Debug("GRPC Unload() received")

	taskID := string(request.GetTaskId())

	return &pluginrpc.UnloadCollectorResponse{}, cs.proxy.UnloadTask(taskID)
}

func (cs *collectService) Info(ctx context.Context, _ *pluginrpc.InfoRequest) (*pluginrpc.InfoResponse, error) {
	logCollectService.Debug("GRPC Info() received")

	pprofAddr := ""
	if cs.pprofLn != nil {
		pprofAddr = cs.pprofLn.Addr().String()
	}

	return serveInfo(ctx, cs.statsController.RequestStat(), pprofAddr)
}

func (cs *collectService) collectWarnings(stream pluginrpc.Collector_CollectServer, warnings []types.Warning) error {
	protoWarnings := make([]*pluginrpc.Warning, 0, len(warnings))

	for _, warn := range warnings {
		protoWarnings = append(protoWarnings, toGRPCWarning(warn))
	}

	if len(warnings) != 0 {
		err := stream.Send(&pluginrpc.CollectResponse{
			Warnings: protoWarnings,
		})
		if err != nil {
			logControlService.WithError(err).Error("can't send warnings chunk over GRPC")
			return err
		}

		logControlService.WithField("len", len(protoWarnings)).Debug("warnings chunk has been sent to snap")
	}

	return nil
}

func (cs *collectService) collectMetrics(stream pluginrpc.Collector_CollectServer, pluginMts []*types.Metric) error {
	protoMts := make([]*pluginrpc.Metric, 0, maxCollectChunkSize)
	for i, pluginMt := range pluginMts {
		protoMt, err := toGRPCMetric(pluginMt)
		if err != nil {
			logCollectService.WithError(err).WithField("metric", pluginMt.Namespace).Errorf("can't send metric over GRPC")
		} else {
			protoMts = append(protoMts, protoMt)
		}

		if len(protoMts) == maxCollectChunkSize || i == len(pluginMts)-1 {
			err = stream.Send(&pluginrpc.CollectResponse{
				MetricSet: protoMts,
			})
			if err != nil {
				logCollectService.WithError(err).Error("can't send metrics chunk over GRPC")
				return err
			}

			logCollectService.WithField("len", len(protoMts)).Debug("metrics chunk has been sent to snap")
			protoMts = make([]*pluginrpc.Metric, 0, len(pluginMts))
		}
	}

	return nil
}
