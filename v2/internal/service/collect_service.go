package service

import (
	"context"
	"fmt"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/pluginrpc"
)

const (
	maxCollectChunkSize = 100
)

var logCollectService = log.WithField("service", "Collect")

type collectService struct {
	proxy CollectorProxy
}

func newCollectService(proxy CollectorProxy) pluginrpc.CollectorServer {
	return &collectService{
		proxy: proxy,
	}
}

func (cs *collectService) Collect(request *pluginrpc.CollectRequest, stream pluginrpc.Collector_CollectServer) error {
	taskID := request.GetTaskId()
	logF := logCollectService.WithField("task-id", taskID)

	logF.Debug("GRPC Collect() received")
	defer logF.Debug("GRPC Collect() completed")

	chunksCh := cs.proxy.RequestCollect(taskID)

	for chunk := range chunksCh {
		err := cs.sendWarnings(stream, chunk.Warnings)
		if err != nil {
			return fmt.Errorf("can't send all warnings to snap: %v", err)
		}

		if chunk.Err != nil {
			return fmt.Errorf("plugin is not able to collect metrics: %s", chunk.Err)
		}

		err = cs.sendMetrics(stream, chunk.Metrics)
		if err != nil {
			return fmt.Errorf("can't send all metrics to snap: %v", err)
		}
	}

	return nil
}

func (cs *collectService) Load(ctx context.Context, request *pluginrpc.LoadCollectorRequest) (*pluginrpc.LoadCollectorResponse, error) {
	taskID := request.GetTaskId()
	logF := logCollectService.WithField("task-id", taskID)

	logF.Debug("GRPC Load() received")
	defer logF.Debug("GRPC Load() completed")

	jsonConfig := request.GetJsonConfig()
	metrics := request.GetMetricSelectors()

	return &pluginrpc.LoadCollectorResponse{}, cs.proxy.LoadTask(taskID, jsonConfig, metrics)
}

func (cs *collectService) Unload(ctx context.Context, request *pluginrpc.UnloadCollectorRequest) (*pluginrpc.UnloadCollectorResponse, error) {
	taskID := request.GetTaskId()
	logF := logCollectService.WithField("task-id", taskID)

	logF.Debug("GRPC Unload() received")
	defer logF.Debug("GRPC Unload() completed")

	return &pluginrpc.UnloadCollectorResponse{}, cs.proxy.UnloadTask(taskID)
}

func (cs *collectService) Info(ctx context.Context, request *pluginrpc.InfoRequest) (*pluginrpc.InfoResponse, error) {
	taskID := request.GetTaskId()
	logF := logCollectService.WithField("task-id", taskID)

	logF.Debug("GRPC Info() received")
	defer logF.Debug("GRPC Info() completed")

	cInfo, err := cs.proxy.CustomInfo(taskID)
	if err != nil {
		return nil, err
	}

	return &pluginrpc.InfoResponse{Info: cInfo}, nil
}

func (cs *collectService) sendWarnings(stream pluginrpc.Collector_CollectServer, warnings []types.Warning) error {
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

func (cs *collectService) sendMetrics(stream pluginrpc.Collector_CollectServer, pluginMts []*types.Metric) error {
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
