// +build medium

package runner

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/rpc"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
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

func (s *SuiteT) sendLoad(taskID int, jsonConfig string, mtsSelector []string) (*rpc.LoadResponse, error) {
	response, err := s.collectorClient.Load(context.Background(), &rpc.LoadRequest{
		TaskId:         int32(taskID),
		JsonConfig:     jsonConfig, // `{"address": {"ip": "127.0.2.3", "port": "12343"}}`,
		MetricSelector: mtsSelector,
	})
	return response, err
}

func (s *SuiteT) sendUnload(taskID int) (*rpc.UnloadResponse, error) {
	response, err := s.collectorClient.Unload(context.Background(), &rpc.UnloadRequest{
		TaskId: int32(taskID),
	})
	return response, err
}

func (s *SuiteT) sendCollect(taskID int) (*rpc.CollectResponse, error) {
	stream, err := s.collectorClient.Collect(context.Background(), &rpc.CollectRequest{
		TaskId: int32(taskID),
	})

	if err != nil {
		return &rpc.CollectResponse{}, err
	}

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &rpc.CollectResponse{}, err
		}
	}

	return &rpc.CollectResponse{}, nil
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

		Convey("Client is able to send load request", func() {
			// Act
			loadResponse, loadErr := s.sendLoad(1, `{"address": {"ip": "127.0.2.3", "port": "12343"}}`, []string{""})

			// Assert
			So(loadErr, ShouldBeNil)
			So(loadResponse, ShouldNotBeNil)
		})

		Convey("Client is able to send collect request", func() {
			// Act
			collectResponse, collectErr := s.sendCollect(1)

			// Assert
			So(collectErr, ShouldBeNil)
			So(collectResponse, ShouldNotBeNil)
			So(simpleCollector.collectCalls, ShouldEqual, 1)
		})

		Convey("Client is able to send several collect request once after another", func() {
			// Act
			for i := 0; i < 5; i++ {
				_, _ = s.sendCollect(1)
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
	collectCalls int
}

func (c *longRunningCollector) Collect(ctx plugin.Context) error {
	c.collectCalls++
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
			_, _ = s.sendLoad(1, "{}", []string{})
			_, err := s.sendCollect(1)
			errCh <- err
		}()

		Convey("Client is able to send kill request and receive no-error response", func() {
			// Act
			time.Sleep(1 * time.Second) // wait for load and blocking collect to be executed from goroutine
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

			// Assert that Collect was called
			So(longRunningCollector.collectCalls, ShouldEqual, 1)
		})
	})
}

/*****************************************************************************/

type simpleConfiguredCollector struct {
	collectCalls int
	loadCalls    int
}

func (c *simpleConfiguredCollector) Load(ctx plugin.Context) error {
	c.loadCalls++
	return nil
}

func (c *simpleConfiguredCollector) Collect(ctx plugin.Context) error {
	c.collectCalls++
	return nil
}

func TestCollectorWithConfiguration(t *testing.T) {

}
