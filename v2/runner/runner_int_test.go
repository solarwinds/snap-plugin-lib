// +build medium

package runner

import (
	"context"
	"io"
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

const expectedGracefulShutdownTimeout = 2 * time.Second
const expectedForceShutdownTimeout = 2*time.Second + rpc.GRPCGracefulStopTimeout

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
	stream, err := s.collectorClient.Collect(context.Background(), &rpc.CollectRequest{
		TaskId: 1,
	})

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &rpc.CollectResponse{}, err
		}
	}

	return &rpc.CollectResponse{}, err
}

/*****************************************************************************/

func TestMedium(t *testing.T) {
	suite.Run(t, new(SuiteT))
}

/*****************************************************************************/

type simpleCollector struct {
	collectCalls int
}

func (sc *simpleCollector) Collect(ctx plugin.Context) error {
	sc.collectCalls++
	return nil
}

func (s *SuiteT) TestSimpleCollector() {
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
			case <-time.After(expectedGracefulShutdownTimeout):
				s.T().Fatal("plugin should have been ended")
			}
		})
	})
}

/*****************************************************************************/

type longRunningCollector struct {
}

func (c *longRunningCollector) Collect(ctx plugin.Context) error {
	time.Sleep(1 * time.Minute)
	return nil
}

func (s *SuiteT) TestKillLongRunningCollector() {
	// Arrange
	longRunningCollector := &longRunningCollector{}
	s.startCollector(longRunningCollector)
	s.startClient()

	errCh := make(chan error, 1)

	Convey("Validate ability to kill collector in case processing takes too much time", s.T(), func() {

		// Act - collect is processing for 1 minute, but kill comes right after request. Should unblock after 10s with error.
		go func() {
			_, err := s.sendCollect()
			errCh <- err
		}()

		Convey("Client is able to send kill request and receive no-error response", func() {
			// Act
			killResponse, killErr := s.sendKill()

			// Assert (kill response)
			So(killErr, ShouldBeNil)
			So(killResponse, ShouldNotBeNil)

			// Assert (plugin has stopped working)
			select {
			case <-errCh:
				// ok
			case <-time.After(expectedForceShutdownTimeout):
				s.T().Fatal("plugin should have been ended")
			}
		})
	})
}
