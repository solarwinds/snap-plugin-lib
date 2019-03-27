// +build medium

package runner

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/pluginrpc"
	"github.com/librato/snap-plugin-lib-go/v2/internal/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

/*****************************************************************************/

const expectedGracefulShutdownTimeout = 2 * time.Second
const expectedForceShutdownTimeout = 2*time.Second + pluginrpc.GRPCGracefulStopTimeout

/*****************************************************************************/

type SuiteT struct {
	suite.Suite

	// grpc server side (plugin)
	startedCollector plugin.Collector
	endCh            chan bool

	// grpc client side (snap)
	grpcConnection  *grpc.ClientConn
	controlClient   pluginrpc.ControllerClient
	collectorClient pluginrpc.CollectorClient
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
		pluginrpc.StartGRPCController(contextManager)
		s.endCh <- true
	}()
}

func (s *SuiteT) startClient() {
	s.grpcConnection, _ = grpc.Dial("localhost:56789", grpc.WithInsecure())

	s.collectorClient = pluginrpc.NewCollectorClient(s.grpcConnection)
	s.controlClient = pluginrpc.NewControllerClient(s.grpcConnection)
}

func (s *SuiteT) sendPing() (*pluginrpc.PingResponse, error) {
	response, err := s.controlClient.Ping(context.Background(), &pluginrpc.PingRequest{})
	return response, err
}

func (s *SuiteT) sendKill() (*pluginrpc.KillResponse, error) {
	response, err := s.controlClient.Kill(context.Background(), &pluginrpc.KillRequest{})
	return response, err
}

func (s *SuiteT) sendLoad(taskID int, configJSON []byte, selectors []string) (*pluginrpc.LoadResponse, error) {
	response, err := s.collectorClient.Load(context.Background(), &pluginrpc.LoadRequest{
		TaskId:          int32(taskID),
		JsonConfig:      configJSON,
		MetricSelectors: selectors,
	})
	return response, err
}

func (s *SuiteT) sendUnload(taskID int) (*pluginrpc.UnloadResponse, error) {
	response, err := s.collectorClient.Unload(context.Background(), &pluginrpc.UnloadRequest{
		TaskId: int32(taskID),
	})
	return response, err
}

func (s *SuiteT) sendCollect(taskID int) (*pluginrpc.CollectResponse, error) {
	stream, err := s.collectorClient.Collect(context.Background(), &pluginrpc.CollectRequest{
		TaskId: int32(taskID),
	})
	if err != nil {
		return nil, err
	}

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &pluginrpc.CollectResponse{}, err
		}
	}

	return &pluginrpc.CollectResponse{}, err
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
	s.T().Skip()

	// Arrange
	jsonConfig := []byte(`{
		"address": {
			"ip": "127.0.2.3", 
			"port": "12343"
		}
	}`)

	mtsSelector := []string{
		"/plugin/metric1",
		"plugin/metric2",
		"/plugin/metric3",
		"/plugin/group1/metric4",
	}

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
			loadResponse, loadErr := s.sendLoad(1, jsonConfig, mtsSelector)

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
	s.T().Skip()
	// Arrange
	jsonConfig := []byte(`{}`)
	var mtsSelector []string

	longRunningCollector := &longRunningCollector{}
	s.startCollector(longRunningCollector)
	s.startClient()

	errCh := make(chan error, 1)

	Convey("Validate ability to kill collector in case processing takes too much time", s.T(), func() {

		// Act - collect is processing for 1 minute, but kill comes right after request. Should unblock after 10s with error.
		go func() {
			_, _ = s.sendLoad(1, jsonConfig, mtsSelector)
			_, err := s.sendCollect(1)
			errCh <- err
		}()

		Convey("Client is able to send kill request and receive no-error response", func() {
			// Act
			time.Sleep(2 * time.Second) // Delay needed to be sure that sendLoad() and sendCollect() in goroutine above were requested
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

type configurableCollector struct {
	t *testing.T
}

type storedObj struct {
	count int
}

func (cc *configurableCollector) Load(ctx plugin.Context) error {
	// Arrange - create configuration objects
	ctx.Store("obj1", &storedObj{count: 10})
	ctx.Store("obj2", &storedObj{count: -14})
	return nil
}

func (cc *configurableCollector) Unload(ctx plugin.Context) error {
	return nil
}

func (cc *configurableCollector) Collect(ctx plugin.Context) error {
	Convey("Validate collector can access objects defined during Load() execution", cc.t, func() {
		// Act
		obj1, ok1 := ctx.Load("obj1")
		obj2, ok2 := ctx.Load("obj2")

		// Assert
		So(ok1, ShouldBeTrue)
		So(obj1, ShouldHaveSameTypeAs, &storedObj{})
		So(obj1.(*storedObj).count, ShouldEqual, 10)

		So(ok2, ShouldBeTrue)
		So(obj2, ShouldHaveSameTypeAs, &storedObj{})
		So(obj2.(*storedObj).count, ShouldEqual, -14)
	})

	Convey("Validate collector can access configuration fields", cc.t, func() {
		// Act
		ip, okIp := ctx.Config("address.ip")
		port, okPort := ctx.Config("address.port")
		user, okUser := ctx.Config("user")

		// Assert
		So(okIp, ShouldBeTrue)
		So(okPort, ShouldBeTrue)
		So(okUser, ShouldBeTrue)

		So(ip, ShouldEqual, "127.0.2.3")
		So(port, ShouldEqual, "12343")
		So(user, ShouldEqual, "admin")
	})

	return nil
}

func (s *SuiteT) TestConfigurableCollector() {
	// Arrange
	jsonConfig := []byte(`{
		"address": {
			"ip": "127.0.2.3", 
			"port": "12343"
		},
		"user": "admin"
	}`)

	var mtsSelector []string

	configurableCollector := &configurableCollector{t: s.T()}
	s.startCollector(configurableCollector)
	s.startClient()

	Convey("", s.T(), func() {
		// Act
		_, _ = s.sendLoad(1, jsonConfig, mtsSelector)
		_, _ = s.sendCollect(1)
		_, _ = s.sendCollect(1)

		time.Sleep(2 * time.Second)
		_, _ = s.sendKill()

		// Assert is handled within configurableCollector.Collect() method
	})
}
