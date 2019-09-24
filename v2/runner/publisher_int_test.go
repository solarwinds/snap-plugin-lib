// +build medium

package runner

import (
	"context"
	"fmt"
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
		return fmt.Errorf("At least one ping wasn't sent properly")
	}

	return nil
}

func (s *PublisherMediumSuite) sendKills() error {
	_, errC := s.collectorControlClient.Kill(context.Background(), &pluginrpc.KillRequest{})
	_, errP := s.publisherControlClient.Kill(context.Background(), &pluginrpc.KillRequest{})

	if errC != nil || errP != nil {
		return fmt.Errorf("At least one kill wasn't sent properly")
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

func (s *PublisherMediumSuite) requestCollectPublishCycle(collectTaskID, publishTaskID string) {

}

///////////////////////////////////////////////////////////////////////////////

func TestPublisherMedium(t *testing.T) {
	suite.Run(t, new(PublisherMediumSuite))
}

///////////////////////////////////////////////////////////////////////////////

type simpleCollector2 struct{}

func (s simpleCollector2) Collect(ctx plugin.CollectContext) error {
	log.Trace("Collect")

	_ = ctx.AddMetric("/example/group1/metric1", 11)
	_ = ctx.AddMetric("/example/group1/metric2", 12)
	_ = ctx.AddMetric("/example/group1/metric3", 13)
	_ = ctx.AddMetric("/example/group1/metric4", 14)
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

	Convey("", p.t, func() {

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

		_, err := s.sendCollectorLoad("task-collector-1", []byte("{}"), []string{})
		So(err, ShouldBeNil)

		_, err = s.sendPublisherLoad("task-publisher-1", []byte("{}"))
		So(err, ShouldBeNil)

		err = s.sendPings()
		So(err, ShouldBeNil)

		// do something
		time.Sleep(1 * time.Second)
		// do something

		err = s.sendPings()
		So(err, ShouldBeNil)

		_, err = s.sendCollectorUnload("task-collector-1")
		So(err, ShouldBeNil)

		_, err = s.sendPublisherUnload("task-publisher-1")
		So(err, ShouldBeNil)

		err = s.sendPings()
		So(err, ShouldBeNil)

		err = s.sendKills()
		So(err, ShouldBeNil)

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
}
