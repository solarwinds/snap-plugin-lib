/*
 Copyright (c) 2020 SolarWinds Worldwide, LLC

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
	"io"

	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/publisher/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/log"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/pluginrpc"
	"github.com/sirupsen/logrus"
)

type publishingService struct {
	proxy proxy.Publisher
	ctx   context.Context
}

func newPublishingService(ctx context.Context, proxy proxy.Publisher) pluginrpc.PublisherServer {
	return &publishingService{
		proxy: proxy,
		ctx:   ctx,
	}
}

func (ps *publishingService) Publish(stream pluginrpc.Publisher_PublishServer) error {
	logF := ps.logger()
	logF.Debug("GRPC Publish() received")

	id := ""
	mts := []*types.Metric{}
	response := &pluginrpc.PublishResponse{}

	for {
		publishPartialReq, err := stream.Recv()
		if err != nil {
			if err == io.EOF { // OK, expected end of stream
				break
			}

			return fmt.Errorf("failure when reading from publish stream: %s", err.Error())
		}

		logF.WithField("length", len(publishPartialReq.MetricSet)).Debug("Metrics chunk received from snap")

		id = publishPartialReq.TaskId

		for _, protoMt := range publishPartialReq.MetricSet {
			mt, err := fromGRPCMetric(protoMt)
			if err != nil {
				logF.WithError(err).Error("can't read metric from GRPC stream")
				continue
			}
			mts = append(mts, &mt)
		}
	}

	if len(mts) != 0 {
		logF.WithField("length", len(mts)).Debug("metric will be published")

		status := ps.proxy.RequestPublish(id, mts)

		protoWarnings := make([]*pluginrpc.Warning, 0, len(status.Warnings))
		for _, w := range status.Warnings {
			protoWarnings = append(protoWarnings, toGRPCWarning(w))
		}

		response.Warnings = protoWarnings

		if status.Error != nil {
			_ = stream.SendAndClose(response) // Ignore potential error from stream, since publish error is of higher importance.
			return status.Error
		}
	} else {
		logF.Info("nothing to publish, request will be ignored")
	}

	return stream.SendAndClose(response)
}

func (ps *publishingService) Load(ctx context.Context, request *pluginrpc.LoadPublisherRequest) (*pluginrpc.LoadPublisherResponse, error) {
	ps.logger().Debug("GRPC Load() received")

	taskID := string(request.GetTaskId())
	jsonConfig := request.GetJsonConfig()

	return &pluginrpc.LoadPublisherResponse{}, ps.proxy.LoadTask(taskID, jsonConfig)
}

func (ps *publishingService) Unload(ctx context.Context, request *pluginrpc.UnloadPublisherRequest) (*pluginrpc.UnloadPublisherResponse, error) {
	ps.logger().Debug("GRPC Unload() received")

	taskID := string(request.GetTaskId())

	return &pluginrpc.UnloadPublisherResponse{}, ps.proxy.UnloadTask(taskID)
}

func (ps *publishingService) Info(ctx context.Context, request *pluginrpc.InfoRequest) (*pluginrpc.InfoResponse, error) {
	ps.logger().Debug("GRPC Info() received")

	taskID := request.GetTaskId()

	cInfo, err := ps.proxy.CustomInfo(taskID)
	if err != nil {
		return nil, err
	}

	return &pluginrpc.InfoResponse{Info: cInfo}, nil
}

func (ps *publishingService) logger() logrus.FieldLogger {
	return log.WithCtx(ps.ctx).WithFields(moduleFields).WithField("service", "Publish")
}
