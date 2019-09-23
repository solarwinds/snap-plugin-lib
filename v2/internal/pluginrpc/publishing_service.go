package pluginrpc

import (
	"fmt"
	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/publisher/proxy"
	"io"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
)

const (
	maxPublishChunkSize = 100
)

type publishingService struct {
	proxy proxy.Publisher
}

func newPublishingService(proxy proxy.Publisher) PublisherServer {
	return &publishingService{
		proxy: proxy,
	}
}

func (ps *publishingService) Publish(stream Publisher_PublishServer) error {
	logF := log.WithField("function", "Publish")

	id := ""
	mts := []*types.Metric{}

	logF.Trace("GRPC Publish() received")

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

		err := ps.proxy.RequestPublish(id, mts)
		if err != nil {
			_ = stream.SendAndClose(&PublishResponse{}) // ignore potential error from stream, since publish error is of higher importance
			return err
		}
	} else {
		logF.Info("nothing to publish, request will be ignored")
	}

	return stream.SendAndClose(&PublishResponse{})
}
