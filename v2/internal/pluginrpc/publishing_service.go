package pluginrpc

import (
	"fmt"
	"io"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
)

const (
	maxPublishChunkSize = 100
)

type publishingService struct {
}

func newPublishingService() PublisherServer {
	return &publishingService{}
}

func (ps *publishingService) Publish(stream Publisher_PublishServer) error {
	logF := log.WithField("function", "Publish")
	mts := []types.Metric{}

	logF.Trace("GRPC Publish() received")

	for {
		protoMts, err := stream.Recv()
		if err != nil {
			if err == io.EOF { // OK, expected end of stream
				break
			}

			return fmt.Errorf("failure when reading from publish stream: %s", err.Error())
		}

		logF.WithField("length", len(protoMts.MetricSet)).Debug("Metrics chunk received from snap")

		for _, protoMt := range protoMts.MetricSet {
			mt, err := fromGRPCMetric(protoMt)
			if err != nil {
				logF.WithError(err).Error("can't read metric from GRPC stream")
				continue
			}
			mts = append(mts, mt)
		}
	}

	if len(mts) != 0 {
		logF.WithField("length", len(mts)).Debug("metric will be published")

		// todo: publish everything
	} else {
		logF.Info("nothing to publish, request will be ignored")
	}

	reply := &PublishResponse{}
	return stream.SendAndClose(reply)
}
