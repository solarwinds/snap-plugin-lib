package pluginrpc

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/librato/snap-plugin-lib-go/v2/internals/plugins/common/stats"
	"github.com/librato/snap-plugin-lib-go/v2/internals/plugins/publisher/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internals/util/types"
)

var logPublishService = log.WithField("service", "Publish")

type publishingService struct {
	proxy           proxy.Publisher
	statsController stats.Controller
	pprofLn         net.Listener
}

func newPublishingService(proxy proxy.Publisher, statsController stats.Controller, pprofLn net.Listener) PublisherServer {
	return &publishingService{
		proxy:           proxy,
		statsController: statsController,
		pprofLn:         pprofLn,
	}
}

func (ps *publishingService) Publish(stream Publisher_PublishServer) error {
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

		err := ps.proxy.RequestPublish(id, mts)
		if err != nil {
			_ = stream.SendAndClose(&PublishResponse{}) // ignore potential error from stream, since publish error is of higher importance
			return err
		}
	} else {
		logPublishService.Info("nothing to publish, request will be ignored")
	}

	return stream.SendAndClose(&PublishResponse{})
}

func (ps *publishingService) Load(ctx context.Context, request *LoadPublisherRequest) (*LoadPublisherResponse, error) {
	logPublishService.Debug("GRPC Load() received")

	taskID := string(request.GetTaskId())
	jsonConfig := request.GetJsonConfig()

	return &LoadPublisherResponse{}, ps.proxy.LoadTask(taskID, jsonConfig)
}

func (ps *publishingService) Unload(ctx context.Context, request *UnloadPublisherRequest) (*UnloadPublisherResponse, error) {
	logPublishService.Debug("GRPC Unload() received")

	taskID := string(request.GetTaskId())

	return &UnloadPublisherResponse{}, ps.proxy.UnloadTask(taskID)
}

func (ps *publishingService) Info(ctx context.Context, _ *InfoRequest) (*InfoResponse, error) {
	logCollectService.Debug("GRPC Info() received")

	pprofAddr := ""
	if ps.pprofLn != nil {
		pprofAddr = ps.pprofLn.Addr().String()
	}

	return serveInfo(ctx, ps.statsController.RequestStat(), pprofAddr)
}
