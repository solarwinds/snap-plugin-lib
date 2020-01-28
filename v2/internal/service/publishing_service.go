package service

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/stats"
	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/publisher/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/pluginrpc"
)

var logPublishService = log.WithField("service", "Publish")

type publishingService struct {
	proxy           proxy.Publisher
	statsController stats.Controller
	pprofLn         net.Listener
}

func newPublishingService(proxy proxy.Publisher, statsController stats.Controller, pprofLn net.Listener) pluginrpc.PublisherServer {
	return &publishingService{
		proxy:           proxy,
		statsController: statsController,
		pprofLn:         pprofLn,
	}
}

func (ps *publishingService) Publish(stream pluginrpc.Publisher_PublishServer) error {
	logPublishService.Debug("GRPC Publish() received")

	id := ""
	mts := []*types.Metric{}

	for {
		publishPartialReq, err := stream.Recv()
		if err != nil {
			if err == io.EOF { // OK, expected end of stream
				break
			}

			return fmt.Errorf("failure when reading from publish stream: %s", err.Error())
		}

		logPublishService.WithField("length", len(publishPartialReq.MetricSet)).Debug("Metrics chunk received from snap")

		id = publishPartialReq.TaskId

		for _, protoMt := range publishPartialReq.MetricSet {
			mt, err := fromGRPCMetric(protoMt)
			if err != nil {
				logPublishService.WithError(err).Error("can't read metric from GRPC stream")
				continue
			}
			mts = append(mts, &mt)
		}
	}

	if len(mts) != 0 {
		logPublishService.WithField("length", len(mts)).Debug("metric will be published")

		status := ps.proxy.RequestPublish(id, mts)
		if status.Error != nil {
			protoWarnings := make([]*pluginrpc.Warning, 0, len(status.Warnings))
			for _, w := range status.Warnings {
				protoWarnings = append(protoWarnings, toGRPCWarning(w))
			}

			// Ignore potential error from stream, since publish error is of higher importance.
			_ = stream.SendAndClose(&pluginrpc.PublishResponse{
				// Warnings are not chunked as in the case of collector. That means it's possible
				// to overflow allowed GRPC Message size (~4MB) for a large number of warnings raised in Publish.
				// Fix would require to change Publish rpc method type to in-out-stream (from in-stream).
				Warnings: protoWarnings,
			})
			return status.Error
		}
	} else {
		logPublishService.Info("nothing to publish, request will be ignored")
	}

	return stream.SendAndClose(&pluginrpc.PublishResponse{})
}

func (ps *publishingService) Load(ctx context.Context, request *pluginrpc.LoadPublisherRequest) (*pluginrpc.LoadPublisherResponse, error) {
	logPublishService.Debug("GRPC Load() received")

	taskID := string(request.GetTaskId())
	jsonConfig := request.GetJsonConfig()

	return &pluginrpc.LoadPublisherResponse{}, ps.proxy.LoadTask(taskID, jsonConfig)
}

func (ps *publishingService) Unload(ctx context.Context, request *pluginrpc.UnloadPublisherRequest) (*pluginrpc.UnloadPublisherResponse, error) {
	logPublishService.Debug("GRPC Unload() received")

	taskID := string(request.GetTaskId())

	return &pluginrpc.UnloadPublisherResponse{}, ps.proxy.UnloadTask(taskID)
}

func (ps *publishingService) Info(ctx context.Context, _ *pluginrpc.InfoRequest) (*pluginrpc.InfoResponse, error) {
	logCollectService.Debug("GRPC Info() received")

	pprofAddr := ""
	if ps.pprofLn != nil {
		pprofAddr = ps.pprofLn.Addr().String()
	}

	return serveInfo(ctx, ps.statsController.RequestStat(), pprofAddr)
}
