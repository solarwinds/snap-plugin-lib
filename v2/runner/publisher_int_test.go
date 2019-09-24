// +build medium

package runner

import (
	"context"
	"github.com/librato/snap-plugin-lib-go/v2/internal/pluginrpc"
	collProxy "github.com/librato/snap-plugin-lib-go/v2/internal/plugins/collector/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/collector/stats"
	pubProxy "github.com/librato/snap-plugin-lib-go/v2/internal/plugins/publisher/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"net"
	"testing"
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

func (s *PublisherMediumSuite) sendPing() (*pluginrpc.PingResponse, error) {
	response, err := s.publisherControlClient.Ping(context.Background(), &pluginrpc.PingRequest{})
	return response, err
}

func (s *PublisherMediumSuite) sendKill() (*pluginrpc.KillResponse, error) {
	response, err := s.publisherControlClient.Kill(context.Background(), &pluginrpc.KillRequest{})
	return response, err
}

// todo: Load, Unload for Publisher

func (s *PublisherMediumSuite) requestCollectPublishCycle(collectTaskID, publishTaskID string) {

}

///////////////////////////////////////////////////////////////////////////////

func TestPublisherMedium(t *testing.T) {
	suite.Run(t, new(PublisherMediumSuite))
}

///////////////////////////////////////////////////////////////////////////////

type simpleCollector2 struct{} // todo: change name

func (s simpleCollector2) Collect(ctx plugin.CollectContext) error {
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

	Convey("", s.T(), func() {

	})
}
