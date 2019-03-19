// +build medium

package runner

import (
	"context"
	"testing"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/rpc"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

/*****************************************************************************/

const ExpectedGracefulShutdownTimeout = 3 * time.Second

//const ExpectedForceShutdownTimeout = 15 * time.Second

/*****************************************************************************/

type simpleCollector struct {
	collectCalls int
}

func (sc *simpleCollector) Collect(ctx plugin.Context) error {
	sc.collectCalls++
	return nil
}

/*****************************************************************************/

type SuiteT struct {
	suite.Suite

	// grpc server side (plugin)
	startedCollector plugin.Collector
	endCh            chan bool

	// grpc client side (snap)
	grpcConnection  *grpc.ClientConn
	controlClient   rpc.ControllerClient
	collectorClient rpc.CollectorClient
}

func (s *SuiteT) SetupSuite() {
	logrus.SetLevel(logrus.TraceLevel)
}

func (s *SuiteT) SetupTest() {
	s.endCh = make(chan bool, 1)
}

func (s *SuiteT) TearDownTest() {
	s.endCh = nil
	s.startedCollector = nil
}

func (s *SuiteT) startCollector(collector plugin.Collector) {
	s.startedCollector = collector
	go func() {
		contextManager := proxy.NewContextManager(collector, "simple_collector", "1.0.0")
		rpc.StartGRPCController(contextManager)
		s.endCh <- true
	}()
}

func (s *SuiteT) startClient() {
	s.grpcConnection, _ = grpc.Dial("localhost:56789", grpc.WithInsecure())

	s.collectorClient = rpc.NewCollectorClient(s.grpcConnection)
	s.controlClient = rpc.NewControllerClient(s.grpcConnection)
}

func (s *SuiteT) sendPing() (*rpc.PingResponse, error) {
	response, err := s.controlClient.Ping(context.Background(), &rpc.PingRequest{})
	return response, err
}

func (s *SuiteT) sendKill() (*rpc.KillResponse, error) {
	response, err := s.controlClient.Kill(context.Background(), &rpc.KillRequest{})
	return response, err
}

func (s *SuiteT) sendLoad() (*rpc.LoadResponse, error) {
	response, err := s.collectorClient.Load(context.Background(), &rpc.LoadRequest{
		TaskId:     1,
		JsonConfig: `{"address": {"ip": "127.0.2.3", "port": "12343"}}`,
		MetricSelector: []string{
			"/plugin/metric1",
			"/plugin/metric2",
			"/plugin/metric3",
			"/plugin/group1/metric4",
		},
	})
	return response, err
}

func (s *SuiteT) sendUnload() (*rpc.UnloadResponse, error) {
	response, err := s.collectorClient.Unload(context.Background(), &rpc.UnloadRequest{
		TaskId: 1,
	})
	return response, err
}

func (s *SuiteT) sendCollect() (*rpc.CollectResponse, error) {
	response, err := s.collectorClient.Collect(context.Background(), &rpc.CollectRequest{
		TaskId: 1,
	})

	_, _ = response.Recv()

	return &rpc.CollectResponse{}, err
}

/*****************************************************************************/

func TestMedium(t *testing.T) {
	suite.Run(t, new(SuiteT))
}

func (s *SuiteT) Test_SimpleCollector() {
	// Arrange
	simpleCollector := &simpleCollector{}
	s.startCollector(simpleCollector)
	s.startClient()

	Convey("Validate ability to connect to simplest collector", s.T(), func() {

		Convey("Client is able to send ping and receive no-error response", func() {
			// Act
			pingResponse, pingErr := s.sendPing()

			// Assert
			So(pingErr, ShouldBeNil)
			So(pingResponse, ShouldNotBeNil)
		})

		Convey("Client is able to send collect request", func() {
			// Act
			collectResponse, collectErr := s.sendCollect()

			// Assert
			So(collectErr, ShouldBeNil)
			So(collectResponse, ShouldNotBeNil)
			So(simpleCollector.collectCalls, ShouldEqual, 1)
		})

		Convey("Client is able to send several collect request once after another", func() {
			// Act
			for i := 0; i < 5; i++ {
				_, _ = s.sendCollect()
			}

			// Assert
			So(simpleCollector.collectCalls, ShouldEqual, 6)
		})

		Convey("Client is able to send kill request and receive no-error response", func() {
			// Act
			killResponse, killErr := s.sendKill()

			// Assert (kill response)
			So(killErr, ShouldBeNil)
			So(killResponse, ShouldNotBeNil)

			// Assert (plugin has stopped working)
			select {
			case <-s.endCh:
			// ok
			case <-time.After(ExpectedGracefulShutdownTimeout):
				s.T().Fatal("plugin should have been ended")
			}
		})
	})
}
