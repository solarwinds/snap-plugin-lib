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

	aggregatedMts := &pluginrpc.CollectResponse{
		MetricSet: []*pluginrpc.Metric{},
	}

	for {
		partialResponse, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		aggregatedMts.MetricSet = append(aggregatedMts.MetricSet, partialResponse.MetricSet...)
	}

	return aggregatedMts, nil
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
	s.T().Skip()

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

/*****************************************************************************/

type metricSendingCollector struct {
	t *testing.T
}

func (msc *metricSendingCollector) DefineMetrics(ctx plugin.CollectorDefinition) error {
	ctx.DefineMetric("/kubernetes/pod/[node]/[namespace]/[pod]/status/phase/Pending", true, "this includes time before being bound to a node, as well as time spent pulling images onto the host")
	ctx.DefineMetric("/kubernetes/pod/[node]/[namespace]/[pod]/status/phase/Running", true, "the pod has been bound to a node and all of the containers have been started")
	ctx.DefineMetric("/kubernetes/pod/[node]/[namespace]/[pod]/status/phase/Succeeded", true, "all containers in the pod have voluntarily terminated with a container exit code of 0, and the system is not going to restart any of these containers")
	ctx.DefineMetric("/kubernetes/pod/[node]/[namespace]/[pod]/status/phase/Failed", true, "all containers in the pod have terminated, and at least one container has terminated in a failure")
	ctx.DefineMetric("/kubernetes/pod/[node]/[namespace]/[pod]/status/phase/Unknown", true, "for some reason the state of the pod could not be obtained, typically due to an error in communicating with the host of the pod")
	ctx.DefineMetric("/kubernetes/pod/[node]/[namespace]/[pod]/status/condition/ready", false, "specifies if the pod is ready to serve requests")
	ctx.DefineMetric("/kubernetes/pod/[node]/[namespace]/[pod]/status/condition/scheduled", false, "status of the scheduling process for the pod")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/status/restarts", true, "number of times the container has been restarted")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/status/ready", true, "specifies whether the container has passed its readiness probe")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/status/waiting", true, "value 1 if container is waiting else value 0")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/status/running", true, "value 1 if container is running else value 0")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/status/terminated", true, "value 1 if container is terminated else value 0")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/requested/cpu/cores", true, "The number of requested cpu cores by a container")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/requested/memory/bytes", true, "The number of requested memory bytes by a container")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/limits/cpu/cores", true, "The number of requested cpu cores by a container")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/limits/memory/bytes", true, "The limit on memory to be used by a container in bytes")
	ctx.DefineMetric("/kubernetes/node/[node]/spec/unschedulable", true, "Whether a node can schedule new pods.")
	ctx.DefineMetric("/kubernetes/node/[node]/status/outofdisk", false, "---")
	ctx.DefineMetric("/kubernetes/node/[node]/status/allocatable/cpu/cores", false, "The CPU resources of a node that are available for scheduling.")
	ctx.DefineMetric("/kubernetes/node/[node]/status/allocatable/memory/bytes", false, "The memory resources of a node that are available for scheduling.")
	ctx.DefineMetric("/kubernetes/node/[node]/status/allocatable/pods", false, "The pod resources of a node that are available for scheduling.")
	ctx.DefineMetric("/kubernetes/node/[node]/status/capacity/cpu/cores", false, "The total CPU resources of the node.")
	ctx.DefineMetric("/kubernetes/node/[node]/status/capacity/memory/bytes", false, "The total memory resources of the node.")
	ctx.DefineMetric("/kubernetes/node/[node]/status/capacity/pods", false, "The total pod resources of the node.")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/metadata/generation", true, "The desired generation sequence number for deployment. If a deployment succeeds should be the same as the observed generation.")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/status/observedgeneration", true, "The generation sequence number after deployment.")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/status/targetedreplicas", true, "Total number of non-terminated pods targeted by this deployment (their labels match the selector).")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/status/availablereplicas", true, "Total number of available pods (ready for at least minReadySeconds) targeted by this deployment.")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/status/unavailablereplicas", true, "Total number of unavailable pods targeted by this deployment.")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/status/updatedreplicas", true, "---")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/status/deploynotfinished", true, "If desired and observed generation are not the same, then either an ongoing deploy or a failed deploy.")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/spec/desiredreplicas", false, "Number of desired pods.")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/spec/paused", false, "---")

	return nil
}

func (msc *metricSendingCollector) Collect(ctx plugin.Context) error {
	Convey("Validate that metrics are filtered acording to metric definitions and filtering", msc.t, func() {
		So(ctx.AddMetric("/kubernetes/pod/node-125/appoptics1/pod-124/status/phase/Running", 1), ShouldBeNil)   // added
		So(ctx.AddMetric("/kubernetes/pod/node-126/appoptics1/pod-124/status/phase/Running", 1), ShouldBeError) // discarded - filtered (node-126 doesn't match filtered rule)
		So(ctx.AddMetric("/kubernetes/pod/node-126/appoptics1/pod-124/status/plase/Running", 1), ShouldBeError) // discarded - no metric "plase" defined

		So(ctx.AddMetric("/kubernetes/container/appoptics1/node-251/pod-34/mycont155/status/ready", 15), ShouldBeNil)   // added
		So(ctx.AddMetric("/kubernetes/container/loggly/node-251/pod-5174/mycont155/status/ready", 21), ShouldBeNil)     // added
		So(ctx.AddMetric("/kubernetes/container/loggly/node-251/pod-5174/mycont155/status", 1), ShouldBeError)          // discarded - no metric status defined
		So(ctx.AddMetric("/kubernetes/container/loggly/node-251/pod-5174/mycont155/status/checking", 1), ShouldBeError) // discarded - no metric status/checking defined

		So(ctx.AddMetric("/kubernetes/node/node-124/status/outofdisk", 1), ShouldBeNil)             // added
		So(ctx.AddMetric("/kubernetes/node/node-124/status/allocatable/cpu/cores", 1), ShouldBeNil) // added

		So(ctx.AddMetric("/kubernetes/deployment/[namespace=appoptics3]/depl-2322/status/targetedreplicas", 10), ShouldBeNil) // added
		So(ctx.AddMetric("/kubernetes/deployment/[namespace=loggly12]/depl-5402/status/availablereplicas", 20), ShouldBeNil)  // added // todo: investigate this scenarios (why passed, next not)
		//So(ctx.AddMetric("/kubernetes/deployment/[namespace=papertrail15]/depl-52/status/updatedreplicas", 30), ShouldBeNil)  // todo: Should pass
		So(ctx.AddMetric("/kubernetes/deployment/[name=appoptics3]/depl-2322/status/targetedreplicas", 1), ShouldBeError) // dicarded (name != namespace)
	})
	return nil
}

func (s *SuiteT) TestMetricSendingCollector() {
	logrus.SetLevel(logrus.TraceLevel)

	// Arrange
	jsonConfig := []byte(`{}`)
	mtsSelector := []string{
		"/kubernetes/pod/node-125/*/*/status/*/*",
		"/kubernetes/container/*/*/*/{mycont[0-9]{3,}}/status/*",
		"/kubernetes/node/*/status/**",
		"/kubernetes/deployment/[namespace={appoptics[0-9]+}]/*/status/*",
		"/kubernetes/deployment/{loggly[0-9]+}/*/{.*}/*",
		"/kubernetes/deployment/papertrail15/*/*/*",
	}

	sendingCollector := &metricSendingCollector{t: s.T()}
	s.startCollector(sendingCollector)
	s.startClient()

	Convey("", s.T(), func() {
		// Act
		_, err := s.sendLoad(1, jsonConfig, mtsSelector)
		So(err, ShouldBeNil)

		collMts, err := s.sendCollect(1)
		So(err, ShouldBeNil)
		So(len(collMts.MetricSet), ShouldEqual, 7)

		time.Sleep(2 * time.Second)
		_, err = s.sendKill()
		So(err, ShouldBeNil)

		// Assert is handled within metricSendingCollector.Collect() method
	})
}
