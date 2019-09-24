// +build medium

package runner

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/pluginrpc"
	collProxy "github.com/librato/snap-plugin-lib-go/v2/internal/plugins/collector/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/collector/stats"
	pubProxy "github.com/librato/snap-plugin-lib-go/v2/internal/plugins/publisher/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

///////////////////////////////////////////////////////////////////////////////

type PublisherMediumSuite struct {
	suite.Suite

	// grpc server side (publisher)
	endPublisherCh   chan bool
	endControllerCh  chan bool
	startedPublisher plugin.Publisher
	startedCollector plugin.Collector

	// grpc client side (snap)
	publisherGRPCConnection *grpc.ClientConn
	publisherControlClient  pluginrpc.ControllerClient
	publisherClient         pluginrpc.PublisherClient

	collectorGRPCConnection *grpc.ClientConn
	collectorControlClient  pluginrpc.ControllerClient
	collectorClient         pluginrpc.CollectorClient
}

func (s *PublisherMediumSuite) SetupSuite() {
	logrus.SetLevel(logrus.TraceLevel)
}

func (s *PublisherMediumSuite) SetupTest() {
	s.startedCollector = nil
	s.startedPublisher = nil

	s.endControllerCh = nil
	s.endPublisherCh = nil
}

func (s *PublisherMediumSuite) TearDownTest() {
}

func (s *PublisherMediumSuite) startCollector(collector plugin.Collector) net.Listener {
	var ln net.Listener

	s.startedCollector = collector
	ln, _ = net.Listen("tcp", "127.0.0.1:")

	go func() {
		statsController, _ := stats.NewEmptyController()
		contextManager := collProxy.NewContextManager(collector, statsController)
		pluginrpc.StartCollectorGRPC(contextManager, statsController, ln, nil, 0, 0)
		s.endControllerCh <- true
	}()

	return ln
}

func (s *PublisherMediumSuite) startPublisher(publisher plugin.Publisher) net.Listener {
	var ln net.Listener

	s.startedPublisher = publisher
	ln, _ = net.Listen("tcp", "127.0.0.1:")

	go func() {
		contextManager := pubProxy.NewContextManager(publisher)
		pluginrpc.StartPublisherGRPC(contextManager, ln, 0, 0)
		s.endPublisherCh <- true
	}()

	return ln
}

func (s *PublisherMediumSuite) startCollectorClient(addr string) {
	s.collectorGRPCConnection, _ = grpc.Dial(addr, grpc.WithInsecure())

	s.collectorClient = pluginrpc.NewCollectorClient(s.collectorGRPCConnection)
	s.collectorControlClient = pluginrpc.NewControllerClient(s.collectorGRPCConnection)
}

func (s *PublisherMediumSuite) startPublisherClient(addr string) {
	s.publisherGRPCConnection, _ = grpc.Dial(addr, grpc.WithInsecure())

	s.publisherClient = pluginrpc.NewPublisherClient(s.publisherGRPCConnection)
	s.publisherControlClient = pluginrpc.NewControllerClient(s.publisherGRPCConnection)
}

func (s *PublisherMediumSuite) sendPings() error {
	_, errC := s.collectorControlClient.Ping(context.Background(), &pluginrpc.PingRequest{})
	_, errP := s.publisherControlClient.Ping(context.Background(), &pluginrpc.PingRequest{})

	if errC != nil || errP != nil {
		return fmt.Errorf("at least one ping wasn't sent properly")
	}

	return nil
}

func (s *PublisherMediumSuite) sendKills() error {
	_, errC := s.collectorControlClient.Kill(context.Background(), &pluginrpc.KillRequest{})
	_, errP := s.publisherControlClient.Kill(context.Background(), &pluginrpc.KillRequest{})

	if errC != nil || errP != nil {
		return fmt.Errorf("at least one kill wasn't sent properly")
	}

	return nil
}

func (s *PublisherMediumSuite) sendCollectorLoad(taskID string, configJSON []byte, selectors []string) (*pluginrpc.LoadCollectorResponse, error) {
	response, err := s.collectorClient.Load(context.Background(), &pluginrpc.LoadCollectorRequest{
		TaskId:          taskID,
		JsonConfig:      configJSON,
		MetricSelectors: selectors,
	})
	return response, err
}

func (s *PublisherMediumSuite) sendPublisherLoad(taskID string, configJSON []byte) (*pluginrpc.LoadPublisherResponse, error) {
	response, err := s.publisherClient.Load(context.Background(), &pluginrpc.LoadPublisherRequest{
		TaskId:     taskID,
		JsonConfig: configJSON,
	})
	return response, err
}

func (s *PublisherMediumSuite) sendCollectorUnload(taskID string) (*pluginrpc.UnloadCollectorResponse, error) {
	response, err := s.collectorClient.Unload(context.Background(), &pluginrpc.UnloadCollectorRequest{
		TaskId: taskID,
	})
	return response, err
}

func (s *PublisherMediumSuite) sendPublisherUnload(taskID string) (*pluginrpc.UnloadPublisherResponse, error) {
	response, err := s.publisherClient.Unload(context.Background(), &pluginrpc.UnloadPublisherRequest{
		TaskId: taskID,
	})
	return response, err
}

func (s *PublisherMediumSuite) requestCollectPublishCycle(collectTaskID, publishTaskID string) error {
	mts, err := s.requestCollect(collectTaskID)
	if err != nil {
		return fmt.Errorf("error when requesting collect from plugin: %v", err)
	}

	err = s.requestPublish(publishTaskID, mts)
	if err != nil {
		return fmt.Errorf("error when requesting publish from plugin: %v", err)
	}

	return nil
}

func (s *PublisherMediumSuite) requestCollect(collectTaskID string) ([]*pluginrpc.Metric, error) {
	stream, err := s.collectorClient.Collect(context.Background(), &pluginrpc.CollectRequest{
		TaskId: collectTaskID,
	})
	if err != nil {
		return nil, err
	}

	mts := []*pluginrpc.Metric{}

	for {
		partialResponse, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		mts = append(mts, partialResponse.MetricSet...)
	}

	return mts, nil
}

func (s *PublisherMediumSuite) requestPublish(publishTaskID string, mts []*pluginrpc.Metric) error {
	stream, err := s.publisherClient.Publish(context.Background())
	if err != nil {
		return nil
	}

	// simplified streaming - only one chunk is sent
	err = stream.Send(&pluginrpc.PublishRequest{
		TaskId:    publishTaskID,
		MetricSet: mts,
	})

	_, err = stream.CloseAndRecv()
	return err
}

///////////////////////////////////////////////////////////////////////////////

func TestPublisherMedium(t *testing.T) {
	suite.Run(t, new(PublisherMediumSuite))
}

///////////////////////////////////////////////////////////////////////////////

type simpleCollector2 struct{}

func (s simpleCollector2) Collect(ctx plugin.CollectContext) error {
	log.Trace("Collect")

	_ = ctx.AddMetricWithTags("/example/group1/metric1", 11, map[string]string{"k1": "v1"})
	_ = ctx.AddMetric("/example/group1/metric2", 12)
	_ = ctx.AddMetric("/example/group1/metric3", 13)
	_ = ctx.AddMetricWithTags("/example/group1/metric4", 14, map[string]string{"k2": "v2", "k3": "v3"})
	_ = ctx.AddMetric("/example/group1/metric5", 15)

	_ = ctx.AddMetric("/example/group2/metric1", 21)
	_ = ctx.AddMetric("/example/group2/metric2", 22)

	return nil
}

type simplePublisher struct {
	t *testing.T
}

func (p simplePublisher) Publish(ctx plugin.PublishContext) error {
	log.Trace("Publish")

	Convey("Validate that all collected metrics are accessed from Publish", p.t, func() {
		// Act
		mts := ctx.ListAllMetrics()

		// Assert
		So(ctx.Count(), ShouldEqual, 7)
		So(len(mts), ShouldEqual, 7)

		So(mts[0].NamespaceText(), ShouldEqual, "/example/group1/metric1")
		So(mts[0].Value(), ShouldEqual, 11)
		So(mts[0].Tags(), ShouldResemble, map[string]string{"k1": "v1"})

		So(mts[1].NamespaceText(), ShouldEqual, "/example/group1/metric2")
		So(mts[1].Value(), ShouldEqual, 12)
		So(mts[1].Tags(), ShouldBeNil)

		So(mts[2].NamespaceText(), ShouldEqual, "/example/group1/metric3")
		So(mts[2].Value(), ShouldEqual, 13)
		So(mts[2].Tags(), ShouldBeNil) // todo: nil vs empty map

		So(mts[3].NamespaceText(), ShouldEqual, "/example/group1/metric4")
		So(mts[3].Value(), ShouldEqual, 14)
		So(mts[3].Tags(), ShouldResemble, map[string]string{"k3": "v3", "k2": "v2"})

		So(mts[4].NamespaceText(), ShouldEqual, "/example/group1/metric5")
		So(mts[4].Value(), ShouldEqual, 15)
		So(mts[4].Tags(), ShouldBeNil)

		So(mts[5].NamespaceText(), ShouldEqual, "/example/group2/metric1")
		So(mts[5].Value(), ShouldEqual, 21)
		So(mts[5].Tags(), ShouldBeNil)

		So(mts[6].NamespaceText(), ShouldEqual, "/example/group2/metric2")
		So(mts[6].Value(), ShouldEqual, 22)
		So(mts[6].Tags(), ShouldBeNil)
	})

	return nil
}

func (s *PublisherMediumSuite) TestSimplePublisher() {
	// Arrange
	collector := &simpleCollector2{}
	publisher := &simplePublisher{t: s.T()}

	lnColl := s.startCollector(collector) // collector server (plugin)
	lnPub := s.startPublisher(publisher)  // publisher server (plugin)

	s.startCollectorClient(lnColl.Addr().String()) // collector client (snap)
	s.startPublisherClient(lnPub.Addr().String())  // collector client (snap)

	Convey("Test that publisher can process all metrics produced by collector", s.T(), func() {

		Convey("Loading tasks for collector and publisher", func() {
			_, err := s.sendCollectorLoad("task-collector-1", []byte("{}"), []string{})
			So(err, ShouldBeNil)

			_, err = s.sendPublisherLoad("task-publisher-1", []byte("{}"))
			So(err, ShouldBeNil)
		})

		Convey("Sending pings requests for collector and publisher (1)", func() {
			err := s.sendPings()
			So(err, ShouldBeNil)
		})

		Convey("Publisher can process all metrics gathered by collector", func() {
			err := s.requestCollectPublishCycle("task-collector-1", "task-publisher-1")
			So(err, ShouldBeNil)

			// validation is done in Publish method of publisher
		})

		Convey("Sending pings requests for collector and publisher (2)", func() {
			err := s.sendPings()
			So(err, ShouldBeNil)
		})

		Convey("Unloading tasks from collector and publisher", func() {
			_, err := s.sendCollectorUnload("task-collector-1")
			So(err, ShouldBeNil)

			_, err = s.sendPublisherUnload("task-publisher-1")
			So(err, ShouldBeNil)
		})

		Convey("Sending pings requests for collector and publisher (3)", func() {
			err := s.sendPings()
			So(err, ShouldBeNil)
		})

		Convey("Sending kills requests for collector and publisher", func() {
			err := s.sendKills()
			So(err, ShouldBeNil)
		})

		Convey("Validate that both plugins quit", func() {
			completeCh := make(chan bool, 1)

			go func() {
				for i := 0; i < 2; i++ {
					select {
					case <-s.endControllerCh:
						// ok
					case <-s.endPublisherCh:
						// ok
					case <-time.After(3 * time.Second):
						break
					}
				}

				completeCh <- true
			}()

			select {
			case <-completeCh:
			// ok
			case <-time.After(10 * time.Second):
				s.T().Fatal("plugin should have been ended")
			}
		})
	})
}
