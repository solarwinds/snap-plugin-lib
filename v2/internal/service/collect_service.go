/*
 Copyright (c) 2021 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package service

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/log"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/types"
	"github.com/solarwinds/snap-plugin-lib/v2/pluginrpc"
)

type collectService struct {
	proxy            CollectorProxy
	ctx              context.Context
	collectChunkSize uint
}

func newCollectService(ctx context.Context, proxy CollectorProxy, collectChunkSize uint) pluginrpc.CollectorServer {
	return &collectService{
		ctx:              ctx,
		proxy:            proxy,
		collectChunkSize: collectChunkSize,
	}
}

func (cs *collectService) Collect(request *pluginrpc.CollectRequest, stream pluginrpc.Collector_CollectServer) error {
	taskID := request.GetTaskId()
	logF := cs.logger().WithField("task-id", taskID)

	logF.Debug("GRPC Collect() received")
	defer logF.Debug("GRPC Collect() completed")

	chunksCh := cs.proxy.RequestCollect(taskID)

	for chunk := range chunksCh {
		// try to send metrics first, even if there were errors during Collect or StreamingCollect
		err := cs.sendMetrics(stream, chunk.Metrics)
		if err != nil {
			return fmt.Errorf("can't send all metrics to snap: %v", err)
		}

		err = cs.sendWarnings(stream, chunk.Warnings)
		if err != nil {
			return fmt.Errorf("can't send all warnings to snap: %v", err)
		}

		if chunk.Err != nil {
			return fmt.Errorf("plugin errored while collecting metrics: %s", chunk.Err)
		}
	}

	return nil
}

func (cs *collectService) Load(ctx context.Context, request *pluginrpc.LoadCollectorRequest) (*pluginrpc.LoadCollectorResponse, error) {
	taskID := request.GetTaskId()
	logF := cs.logger().WithField("task-id", taskID)

	logF.Debug("GRPC Load() received")
	defer logF.Debug("GRPC Load() completed")

	jsonConfig := request.GetJsonConfig()
	metrics := request.GetMetricSelectors()

	return &pluginrpc.LoadCollectorResponse{}, cs.proxy.LoadTask(taskID, jsonConfig, metrics)
}

func (cs *collectService) Unload(ctx context.Context, request *pluginrpc.UnloadCollectorRequest) (*pluginrpc.UnloadCollectorResponse, error) {
	taskID := request.GetTaskId()
	logF := cs.logger().WithField("task-id", taskID)

	logF.Debug("GRPC Unload() received")
	defer logF.Debug("GRPC Unload() completed")

	return &pluginrpc.UnloadCollectorResponse{}, cs.proxy.UnloadTask(taskID)
}

func (cs *collectService) Info(ctx context.Context, request *pluginrpc.InfoRequest) (*pluginrpc.InfoResponse, error) {
	taskID := request.GetTaskId()
	logF := cs.logger().WithField("task-id", taskID)

	logF.Debug("GRPC Info() received")
	defer logF.Debug("GRPC Info() completed")

	cInfo, err := cs.proxy.CustomInfo(taskID)
	if err != nil {
		return nil, err
	}

	return &pluginrpc.InfoResponse{Info: cInfo}, nil
}

func (cs *collectService) sendWarnings(stream pluginrpc.Collector_CollectServer, warnings []types.Warning) error {
	logF := cs.logger()
	protoWarnings := make([]*pluginrpc.Warning, 0, len(warnings))

	for _, warn := range warnings {
		protoWarnings = append(protoWarnings, toGRPCWarning(warn))
	}

	if len(warnings) != 0 {
		err := stream.Send(&pluginrpc.CollectResponse{
			Warnings: protoWarnings,
		})
		if err != nil {
			logF.WithError(err).Error("can't send warnings chunk over GRPC")
			return err
		}

		logF.WithField("len", len(protoWarnings)).Debug("warnings chunk has been sent to snap")
	}

	return nil
}

func (cs *collectService) sendMetrics(stream pluginrpc.Collector_CollectServer, pluginMts []*types.Metric) error {
	logF := cs.logger()

	protoMts := make([]*pluginrpc.Metric, 0, cs.collectChunkSize)
	for i, pluginMt := range pluginMts {
		protoMt, err := toGRPCMetric(pluginMt)
		if err != nil {
			logF.WithError(err).WithField("metric", pluginMt.Namespace).Errorf("can't send metric over GRPC")
		} else {
			protoMts = append(protoMts, protoMt)
		}

		if len(protoMts) == int(cs.collectChunkSize) || i == len(pluginMts)-1 {
			err = stream.Send(&pluginrpc.CollectResponse{
				MetricSet: protoMts,
			})
			if err != nil {
				logF.WithError(err).Error("can't send metrics chunk over GRPC")
				return err
			}

			logF.WithField("len", len(protoMts)).Debug("metrics chunk has been sent to snap")
			protoMts = make([]*pluginrpc.Metric, 0, len(pluginMts))
		}
	}

	return nil
}

func (cs *collectService) logger() logrus.FieldLogger {
	return log.WithCtx(cs.ctx).WithFields(moduleFields).WithField("service", "Collect")
}
