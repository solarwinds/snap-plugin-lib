// +build medium

/*
 Copyright (c) 2020 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package runner

import (
	"context"
	"fmt"
	"io"
	"math"
	"net"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/plugins/collector/proxy"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/plugins/common/stats"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/service"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/types"
	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
	"github.com/solarwinds/snap-plugin-lib/v2/pluginrpc"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

/*****************************************************************************/

const expectedGracefulShutdownTimeout = 2 * time.Second
const expectedForceShutdownTimeout = 2*time.Second + service.GRPCGracefulStopTimeout
const expectedUnloadTimeout = 3 * time.Second

/*****************************************************************************/

type SuiteT struct {
	suite.Suite

	// grpc server side (plugin)
	startedCollector          plugin.Collector
	startedStreamingCollector plugin.StreamingCollector
	endCh                     chan bool

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

func (s *SuiteT) startCollector(collector plugin.Collector) net.Listener {
	var ln net.Listener

	s.startedCollector = collector
	ln, _ = net.Listen("tcp", "127.0.0.1:")

	go func() {
		statsController, _ := stats.NewEmptyController()
		contextManager := proxy.NewContextManager(context.Background(), types.NewCollector("test-collector", "1.0.0", collector), statsController)
		service.StartCollectorGRPC(context.Background(), grpc.NewServer(), contextManager, ln, 0, 0)
		s.endCh <- true
	}()

	return ln
}

func (s *SuiteT) startStreamingCollector(collector plugin.StreamingCollector) net.Listener {
	var ln net.Listener

	s.startedStreamingCollector = collector
	ln, _ = net.Listen("tcp", "127.0.0.1:")

	go func() {
		statsController, _ := stats.NewEmptyController()
		contextManager := proxy.NewContextManager(context.Background(), types.NewStreamingCollector("test-collector", "1.0.0", collector), statsController)
		service.StartCollectorGRPC(context.Background(), grpc.NewServer(), contextManager, ln, 0, 0)
		s.endCh <- true
	}()

	return ln
}

func (s *SuiteT) startClient(addr string) {
	s.grpcConnection, _ = grpc.Dial(addr, grpc.WithInsecure())

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

func (s *SuiteT) sendLoad(taskID string, configJSON []byte, selectors []string) (*pluginrpc.LoadCollectorResponse, error) {
	response, err := s.collectorClient.Load(context.Background(), &pluginrpc.LoadCollectorRequest{
		TaskId:          taskID,
		JsonConfig:      configJSON,
		MetricSelectors: selectors,
	})
	return response, err
}

func (s *SuiteT) sendUnload(taskID string) (*pluginrpc.UnloadCollectorResponse, error) {
	response, err := s.collectorClient.Unload(context.Background(), &pluginrpc.UnloadCollectorRequest{
		TaskId: taskID,
	})
	return response, err
}

func (s *SuiteT) sendCollect(taskID string) (*pluginrpc.CollectResponse, error) {
	stream, err := s.collectorClient.Collect(context.Background(), &pluginrpc.CollectRequest{
		TaskId: taskID,
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

func (sc *simpleCollector) Collect(ctx plugin.CollectContext) error {
	sc.collectCalls++
	return nil
}

func (s *SuiteT) TestSimpleCollector() {
	// Arrange
	jsonConfig := []byte(`{
		"address": {
			"ip": "127.0.2.3", 
			"port": "12343"
		}
	}`)

	mtsSelector := []string{
		"/plugin/metric1",
		"/plugin/metric3",
		"/plugin/group1/metric4",
	}

	simpleCollector := &simpleCollector{}
	ln := s.startCollector(simpleCollector)
	s.startClient(ln.Addr().String())

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
			loadResponse, loadErr := s.sendLoad("task-1", jsonConfig, mtsSelector)

			// Assert
			So(loadErr, ShouldBeNil)
			So(loadResponse, ShouldNotBeNil)
		})

		Convey("Client is able to send collect request", func() {
			// Act
			collectResponse, collectErr := s.sendCollect("task-1")

			// Assert
			So(collectErr, ShouldBeNil)
			So(collectResponse, ShouldNotBeNil)
			So(simpleCollector.collectCalls, ShouldEqual, 1)
		})

		Convey("Client is able to send several collect request once after another", func() {
			// Act
			for i := 0; i < 5; i++ {
				_, _ = s.sendCollect("task-1")
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

func (s *SuiteT) TestMisconfiguredCollector() {
	// Arrange
	jsonConfig := []byte(`{}`)
	mtsSelector := []string{
		"plugin/metric2", // ! wrong namespace
	}

	noDefCollector := &simpleCollector{}
	ln := s.startCollector(noDefCollector)
	s.startClient(ln.Addr().String())

	Convey("Validate that load throws error when requested metrics/filters are invalid", s.T(), func() {
		// Act
		_, err := s.sendLoad("task-1", jsonConfig, mtsSelector)
		So(err, ShouldBeError)
	})
}

/*****************************************************************************/

type longRunningCollector struct {
	collectCalls    int
	collectDuration time.Duration
}

func (c *longRunningCollector) Collect(ctx plugin.CollectContext) error {
	c.collectCalls++
	time.Sleep(c.collectDuration)
	return nil
}

func (s *SuiteT) TestKillLongRunningCollector() {
	// Arrange
	jsonConfig := []byte(`{}`)
	var mtsSelector []string

	longRunningCollector := &longRunningCollector{collectDuration: 1 * time.Second}
	ln := s.startCollector(longRunningCollector)
	s.startClient(ln.Addr().String())

	errCh := make(chan error, 1)

	Convey("Validate ability to kill collector in case processing takes too much time", s.T(), func() {

		// Act - collect is processing for 1 minute, but kill comes right after request. Should unblock after 10s with error.
		go func() {
			_, _ = s.sendLoad("task-1", jsonConfig, mtsSelector)
			_, err := s.sendCollect("task-1")
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

func (s *SuiteT) TestRunningCollectorAtTheSameTime() {
	// Arrange
	jsonConfig := []byte(`{}`)
	var mtsSelector []string

	longRunningCollector := &longRunningCollector{collectDuration: 5 * time.Second}
	ln := s.startCollector(longRunningCollector)
	s.startClient(ln.Addr().String())

	errCh := make(chan error, 1)

	Convey("Validate that collector associated with the same id can't be run in more that 1 instance", s.T(), func() {
		const numberOfCollectors = 10
		const numberOfCollectorsWithSameID = 5

		for id := 1; id <= numberOfCollectors; id++ {
			_, _ = s.sendLoad(fmt.Sprintf("task-%d", id), jsonConfig, mtsSelector)

			for i := 1; i <= numberOfCollectorsWithSameID; i++ {
				go func(id int) {
					_, err := s.sendCollect(fmt.Sprintf("task-%d", id))
					errCh <- err
				}(id)
			}
		}

		errorCounter := 0
		for i := 0; i < numberOfCollectorsWithSameID*numberOfCollectors; i++ {
			errRecv := <-errCh
			if errRecv != nil {
				errorCounter++
			}
		}

		So(errorCounter, ShouldEqual, (numberOfCollectors*numberOfCollectorsWithSameID)-numberOfCollectors) // only 1 task from each id should complete without error

		// validate that when collect is completed you can requested it
		_, err := s.sendCollect("task-1")
		So(err, ShouldBeNil)

		time.Sleep(2 * time.Second)
		_, _ = s.sendKill()
	})
}

/*****************************************************************************/

type configurableCollector struct {
	t           *testing.T
	loadCalls   int
	unloadCalls int
}

type storedObj struct {
	count int
}

func (cc *configurableCollector) resetCallCounters() {
	cc.loadCalls = 0
	cc.unloadCalls = 0
}

func (cc *configurableCollector) Load(ctx plugin.Context) error {
	cc.loadCalls++

	// Arrange - create configuration objects
	ctx.Store("obj1", &storedObj{count: 10})
	ctx.Store("obj2", &storedObj{count: -14})

	return nil
}

func (cc *configurableCollector) Unload(ctx plugin.Context) error {
	cc.unloadCalls++

	return nil
}

func (cc *configurableCollector) Collect(ctx plugin.CollectContext) error {
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
		ip, okIp := ctx.ConfigValue("address.ip")
		port, okPort := ctx.ConfigValue("address.port")
		user, okUser := ctx.ConfigValue("user")

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
	ln := s.startCollector(configurableCollector)
	s.startClient(ln.Addr().String())

	Convey("Validate that load and unload works in valid scenarios", s.T(), func() {
		{
			_, err := s.sendLoad("task-1", jsonConfig, mtsSelector)
			So(err, ShouldBeNil)
			So(configurableCollector.loadCalls, ShouldEqual, 1)
		}
		{
			_, err := s.sendCollect("task-1")
			So(err, ShouldBeNil)
		}
		{
			_, err := s.sendCollect("task-1")
			So(err, ShouldBeNil)
		}
		{
			_, err := s.sendUnload("task-1")
			So(err, ShouldBeNil)
			So(configurableCollector.unloadCalls, ShouldEqual, 1)
		}
	})

	Convey("Validate that load and unload works properly in invalid scenarios", s.T(), func() {
		configurableCollector.resetCallCounters()

		{
			_, err := s.sendLoad("task-1", jsonConfig, mtsSelector)
			So(err, ShouldBeNil)
			So(configurableCollector.loadCalls, ShouldEqual, 1)
		}
		{
			_, err := s.sendLoad("task-2", jsonConfig, mtsSelector)
			So(err, ShouldBeNil)
			So(configurableCollector.loadCalls, ShouldEqual, 2)
		}
		{ // Shouldn't accept load of the same task
			_, err := s.sendLoad("task-1", jsonConfig, mtsSelector)
			So(err, ShouldBeError)
			So(configurableCollector.loadCalls, ShouldEqual, 2)
		}
		{ // Shouldn't accept unload of the same task which was loaded
			_, err := s.sendUnload("task-3")
			So(err, ShouldBeError)
			So(configurableCollector.unloadCalls, ShouldEqual, 0)
		}
		{
			_, err := s.sendUnload("task-1")
			So(err, ShouldBeNil)
			So(configurableCollector.unloadCalls, ShouldEqual, 1)
		}
		{ // Shouldn't accept unload of the task that is already unloaded
			_, err := s.sendUnload("task-1")
			So(err, ShouldBeError)
			So(configurableCollector.unloadCalls, ShouldEqual, 1)
		}
		{
			_, err := s.sendUnload("task-2")
			So(err, ShouldBeNil)
			So(configurableCollector.unloadCalls, ShouldEqual, 2)
		}
	})

	time.Sleep(2 * time.Second)
	_, _ = s.sendKill()

	// Assert is handled within configurableCollector.Collect() method
}

/*****************************************************************************/

type kubernetesCollector struct {
	t *testing.T
}

func (kc *kubernetesCollector) PluginDefinition(ctx plugin.CollectorDefinition) error {
	ctx.DefineMetric("/kubernetes/pod/[node]/[namespace]/[pod]/status/phase/Pending", "", true, "this includes time before being bound to a node, as well as time spent pulling images onto the host")
	ctx.DefineMetric("/kubernetes/pod/[node]/[namespace]/[pod]/status/phase/Running", "count", true, "the pod has been bound to a node and all of the containers have been started")
	ctx.DefineMetric("/kubernetes/pod/[node]/[namespace]/[pod]/status/phase/Succeeded", "", true, "all containers in the pod have voluntarily terminated with a container exit code of 0, and the system is not going to restart any of these containers")
	ctx.DefineMetric("/kubernetes/pod/[node]/[namespace]/[pod]/status/phase/Failed", "", true, "all containers in the pod have terminated, and at least one container has terminated in a failure")
	ctx.DefineMetric("/kubernetes/pod/[node]/[namespace]/[pod]/status/phase/Unknown", "", true, "for some reason the state of the pod could not be obtained, typically due to an error in communicating with the host of the pod")
	ctx.DefineMetric("/kubernetes/pod/[node]/[namespace]/[pod]/status/condition/ready", "", false, "specifies if the pod is ready to serve requests")
	ctx.DefineMetric("/kubernetes/pod/[node]/[namespace]/[pod]/status/condition/scheduled", "", false, "status of the scheduling process for the pod")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/status/restarts", "", true, "number of times the container has been restarted")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/status/ready", "boolean", true, "specifies whether the container has passed its readiness probe")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/status/waiting", "", true, "value 1 if container is waiting else value 0")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/status/running", "", true, "value 1 if container is running else value 0")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/status/terminated", "", true, "value 1 if container is terminated else value 0")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/requested/cpu/cores", "", true, "The number of requested cpu cores by a container")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/requested/memory/bytes", "", true, "The number of requested memory bytes by a container")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/limits/cpu/cores", "", true, "The number of requested cpu cores by a container")
	ctx.DefineMetric("/kubernetes/container/[namespace]/[node]/[pod]/[container]/limits/memory/bytes", "", true, "The limit on memory to be used by a container in bytes")
	ctx.DefineMetric("/kubernetes/node/[node]/spec/unschedulable", "", true, "Whether a node can schedule new pods.")
	ctx.DefineMetric("/kubernetes/node/[node]/status/outofdisk", "", false, "---")
	ctx.DefineMetric("/kubernetes/node/[node]/status/allocatable/cpu/cores", "bytes", false, "The CPU resources of a node that are available for scheduling.")
	ctx.DefineMetric("/kubernetes/node/[node]/status/allocatable/memory/bytes", "bytes", false, "The memory resources of a node that are available for scheduling.")
	ctx.DefineMetric("/kubernetes/node/[node]/status/allocatable/pods", "bytes", false, "The pod resources of a node that are available for scheduling.")
	ctx.DefineMetric("/kubernetes/node/[node]/status/capacity/cpu/cores", "", false, "The total CPU resources of the node.")
	ctx.DefineMetric("/kubernetes/node/[node]/status/capacity/memory/bytes", "", false, "The total memory resources of the node.")
	ctx.DefineMetric("/kubernetes/node/[node]/status/capacity/pods", "", false, "The total pod resources of the node.")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/metadata/generation", "", true, "The desired generation sequence number for deployment. If a deployment succeeds should be the same as the observed generation.")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/status/observedgeneration", "", true, "The generation sequence number after deployment.")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/status/targetedreplicas", "", true, "Total number of non-terminated pods targeted by this deployment (their labels match the selector).")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/status/availablereplicas", "", true, "Total number of available pods (ready for at least minReadySeconds) targeted by this deployment.")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/status/unavailablereplicas", "", true, "Total number of unavailable pods targeted by this deployment.")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/status/updatedreplicas", "", true, "---")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/status/deploynotfinished", "", true, "If desired and observed generation are not the same, then either an ongoing deploy or a failed deploy.")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/spec/desiredreplicas", "", false, "Number of desired pods.")
	ctx.DefineMetric("/kubernetes/deployment/[namespace]/[deployment]/spec/paused", "", false, "---")

	ctx.DefineGroup("node", "kubernetes node name")
	ctx.DefineGroup("namespace", "kubernetes namespace")
	ctx.DefineGroup("pod", "kubernetes pod")
	ctx.DefineGroup("container", "kubernetes container")
	ctx.DefineGroup("deployment", "kubernetes deployment")

	return nil
}

func (kc *kubernetesCollector) Collect(ctx plugin.CollectContext) error {
	Convey("Validate that user can obtain proper information about reasonableness to process metrics or metrics groups", kc.t, func() {
		So(ctx.ShouldProcess("/kubernetes/deployment/*/*/spec/paused"), ShouldBeTrue)                  // ok
		So(ctx.ShouldProcess("/kubernetes/deployment/*/*/spec/*"), ShouldBeTrue)                       // ok
		So(ctx.ShouldProcess("/kubernetes/deployment/*/*/*/paused"), ShouldBeTrue)                     // ok
		So(ctx.ShouldProcess("/kubernetes/deployment/*/*/*/*"), ShouldBeTrue)                          // ok
		So(ctx.ShouldProcess("/kubernetes/deployment/*/*/*/*/*"), ShouldBeFalse)                       // false - too much elements
		So(ctx.ShouldProcess("/kubernetes/deployment/*/*/*"), ShouldBeTrue)                            // ok
		So(ctx.ShouldProcess("/kubernetes/deployment/*/*/spec/paused/*"), ShouldBeFalse)               // false - too much elements
		So(ctx.ShouldProcess("/kubernetes/deployment/*/depl-01/spec/paused"), ShouldBeTrue)            // ok
		So(ctx.ShouldProcess("/kubernetes/deployment/papertrail15/*/spec/paused"), ShouldBeTrue)       // ok
		So(ctx.ShouldProcess("/kubernetes/deployment/papertrail16/*/spec/paused"), ShouldBeFalse)      // false - only "papertrail15" is allowed by filters
		So(ctx.ShouldProcess("/kubernetes/deployment/papertrail15/depl-01/spec/paused"), ShouldBeTrue) // ok
		So(ctx.ShouldProcess("/kubernetes/deployment/papertrail15/depl-01/spec"), ShouldBeTrue)        // ok
		So(ctx.ShouldProcess("/kubernetes/deployment/papertrail15/depl-01"), ShouldBeTrue)             // ok
		So(ctx.ShouldProcess("/kubernetes/deployment/papertrail15"), ShouldBeTrue)                     // ok
		So(ctx.ShouldProcess("/kubernetes/deployment/papertrail16"), ShouldBeFalse)                    // false - only "papertrail15" is allowed by filters
		So(ctx.ShouldProcess("/kubernetes/deployment/*"), ShouldBeTrue)                                // ok
		So(ctx.ShouldProcess("/kubernetes/deployment"), ShouldBeTrue)                                  // ok
		So(ctx.ShouldProcess("/kubernetes"), ShouldBeFalse)                                            // false - ns should have at least 2 elements

		// Following checks are adequate to results of AddMetrics from next convey section with some exceptions
		So(ctx.ShouldProcess("/kubernetes/pod/node-125/appoptics1/pod-124/status/phase/Running"), ShouldBeTrue)
		So(ctx.ShouldProcess("/kubernetes/pod/node-126/appoptics1/pod-124/status/phase/Running"), ShouldBeFalse)
		So(ctx.ShouldProcess("/kubernetes/pod/node-126/appoptics1/pod-124/status/plase/Running"), ShouldBeFalse)

		So(ctx.ShouldProcess("/kubernetes/container/appoptics1/node-251/pod-34/mycont155/status/ready"), ShouldBeTrue)
		So(ctx.ShouldProcess("/kubernetes/container/loggly/node-251/pod-5174/mycont155/status/ready"), ShouldBeTrue)
		So(ctx.ShouldProcess("/kubernetes/container/loggly/node-251/pod-5174/mycont155/status"), ShouldBeTrue) // ok, because last element doesn't have to be metrics name
		So(ctx.ShouldProcess("/kubernetes/container/loggly/node-251/pod-5174/mycont155/status/checking"), ShouldBeFalse)

		So(ctx.ShouldProcess("/kubernetes/node/node-124/status/outofdisk"), ShouldBeTrue)
		So(ctx.ShouldProcess("/kubernetes/node/node-124/status/allocatable/cpu/cores"), ShouldBeTrue)

		So(ctx.ShouldProcess("/kubernetes/deployment/[namespace=appoptics3]/depl-2322/status/targetedreplicas"), ShouldBeTrue)
		So(ctx.ShouldProcess("/kubernetes/deployment/[namespace=loggly12]/depl-5402/status/availablereplicas"), ShouldBeTrue)
		So(ctx.ShouldProcess("/kubernetes/deployment/[namespace=papertrail15]/depl-52/status/updatedreplicas"), ShouldBeTrue)
		So(ctx.ShouldProcess("/kubernetes/deployment/[name=appoptics3]/depl-2322/status/targetedreplicas"), ShouldBeFalse)
	})

	Convey("Validate that metrics are filtered according to metric definitions and filtering", kc.t, func() {
		So(ctx.AddMetric("/kubernetes/pod/node-125/appoptics1/pod-124/status/phase/Running", 1), ShouldBeNil)   // added
		So(ctx.AddMetric("/kubernetes/pod/node-126/appoptics1/pod-124/status/phase/Running", 1), ShouldBeNil)   // discarded - filtered (node-126 doesn't match filtered rule)
		So(ctx.AddMetric("/kubernetes/pod/node-126/appoptics1/pod-124/status/plase/Running", 1), ShouldBeError) // discarded - no metric "plase" defined

		So(ctx.AddMetric("/kubernetes/container/appoptics1/node-251/pod-34/mycont155/status/ready", 15), ShouldBeNil)   // added
		So(ctx.AddMetric("/kubernetes/container/loggly/node-251/pod-5174/mycont155/status/ready", 21), ShouldBeNil)     // added
		So(ctx.AddMetric("/kubernetes/container/loggly/node-251/pod-5174/mycont155/status", 1), ShouldBeError)          // discarded - no metric status defined
		So(ctx.AddMetric("/kubernetes/container/loggly/node-251/pod-5174/mycont155/status/checking", 1), ShouldBeError) // discarded - no metric status/checking defined

		So(ctx.AddMetric("/kubernetes/node/node-124/status/outofdisk", 1), ShouldBeNil)             // added
		So(ctx.AddMetric("/kubernetes/node/node-124/status/allocatable/cpu/cores", 1), ShouldBeNil) // added

		So(ctx.AddMetric("/kubernetes/deployment/[namespace=appoptics3]/depl-2322/status/targetedreplicas", 10), ShouldBeNil) // added
		So(ctx.AddMetric("/kubernetes/deployment/[namespace=loggly12]/depl-5402/status/availablereplicas", 20), ShouldBeNil)  // added
		So(ctx.AddMetric("/kubernetes/deployment/[namespace=papertrail15]/depl-52/status/updatedreplicas", 30), ShouldBeNil)  // added
		So(ctx.AddMetric("/kubernetes/deployment/[name=appoptics3]/depl-2322/status/targetedreplicas", 1), ShouldBeError)     // discarded (name != namespace)
	})

	Convey("Validate collector context can provide requested metrics", kc.t, func() {
		// Arrange
		expectedReqMts := []string{
			"/kubernetes/container/*/*/*/{mycont[0-9]{3,}}/status/*",
			"/kubernetes/deployment/[namespace={appoptics[0-9]+}]/*/status/*",
			"/kubernetes/deployment/papertrail15/*/*/*",
			"/kubernetes/deployment/{loggly[0-9]+}/*/{.*}/*",
			"/kubernetes/node/*/status/**",
			"/kubernetes/pod/node-125/*/*/status/*/*",
		}

		// Act
		reqMts := ctx.RequestedMetrics()

		// Assert
		So(reqMts, ShouldResemble, expectedReqMts)
	})

	return nil
}

func (s *SuiteT) TestKubernetesCollector() {
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

	sendingCollector := &kubernetesCollector{t: s.T()}
	ln := s.startCollector(sendingCollector)
	s.startClient(ln.Addr().String())

	Convey("Validate that collector can gather metric based on definition and filter", s.T(), func() {
		// Act
		_, err := s.sendLoad("task-1", jsonConfig, mtsSelector)
		So(err, ShouldBeNil)

		collMts, err := s.sendCollect("task-1")
		So(err, ShouldBeNil)
		So(len(collMts.MetricSet), ShouldEqual, 8)

		time.Sleep(2 * time.Second)
		_, err = s.sendKill()
		So(err, ShouldBeNil)

		// Assert is handled within kubernetesCollector.Collect() method
	})
}

func (s *SuiteT) TestMisconfiguredCollectorWithDefinedMetrics() {
	// Arrange
	jsonConfig := []byte(`{}`)
	mtsSelector := []string{
		"/kubernetes/pod/node-125/*/*/status/*/*",
		"/kubernetes/container/*/*/*/{mycont[0-9]{3,}}/status/*",
		"/kubernetes/node/*/status/**",
		"kubernetes/deployment/[namespace={appoptics[0-9]+}]/*/status/*", // ! wrong
		"/kubernetes/deployment/{loggly[0-9]+}/*/{.*}/*",
		"/kubernetes/deployment/papertrail15/*/*/*",
	}

	noDefCollector := &kubernetesCollector{}
	ln := s.startCollector(noDefCollector)
	s.startClient(ln.Addr().String())

	Convey("Validate that load throws error when requested metrics/filters are invalid", s.T(), func() {
		// Act
		_, err := s.sendLoad("task-1", jsonConfig, mtsSelector)
		So(err, ShouldBeError)
	})
}

/*****************************************************************************/

type noDefinitionCollector struct {
	collectCalls int
	t            *testing.T
}

func (ndc *noDefinitionCollector) Collect(ctx plugin.CollectContext) error {
	ndc.collectCalls++

	Convey("Validate that user can obtain proper information about reasonableness to process metrics or metrics groups", ndc.t, func() {
		// Following checks are adequate to results of AddMetrics from next convey section
		So(ctx.ShouldProcess("/plugin/group1/subgroup1/metric1"), ShouldBeTrue)
		So(ctx.ShouldProcess("/plugin/group2/id12/metric1"), ShouldBeTrue)
		So(ctx.ShouldProcess("/plugin/group3/subgroup3/metric4"), ShouldBeTrue)
		So(ctx.ShouldProcess("/plugin/group3/subgroup3/metric$4"), ShouldBeFalse)
		So(ctx.ShouldProcess("/plugin/group3/subgroup4/metric4"), ShouldBeTrue)
		So(ctx.ShouldProcess("/plugin/group3/subgroup4/sub5/metric6"), ShouldBeTrue)
		So(ctx.ShouldProcess("/plugin/group3/subgroup4/sub()5/metric6"), ShouldBeFalse)
		So(ctx.ShouldProcess("some/plugin/group1/subgroup1/metric1"), ShouldBeFalse)
		So(ctx.ShouldProcess("/plugin/group2/[subgroup2=id12]/metric1"), ShouldBeFalse)
	})

	Convey("Validate that metrics are filtered according to filtering", ndc.t, func() {
		So(ctx.AddMetric("/plugin/group1/subgroup1/metric1", 10,
			plugin.MetricDescription("metric1 custom description")), ShouldBeNil) // added

		So(ctx.AddMetric("/plugin/group2/id12/metric1", 20,
			plugin.MetricUnit("MU"),
			plugin.MetricTimestamp(time.Date(2020, 10, 13, 12, 13, 14, 0, time.UTC))), ShouldBeNil) // added

		So(ctx.AddMetric("/plugin/group3/subgroup3/metric4", 15,
			plugin.MetricTags(map[string]string{"a": "b", "c": "d"})), ShouldBeNil) // added

		So(ctx.AddMetric("/plugin/group3/subgroup3/metric$4", 12), ShouldBeError)       // invalid char used in element
		So(ctx.AddMetric("/plugin/group3/subgroup4/metric4", 15), ShouldBeNil)          // added
		So(ctx.AddMetric("/plugin/group3/subgroup4/sub5/metric6", 13), ShouldBeNil)     // added
		So(ctx.AddMetric("/plugin/group3/subgroup4/sub()5/metric6", 13), ShouldBeError) // invalid char used in element
		So(ctx.AddMetric("some/plugin/group1/subgroup1/metric1", 11), ShouldBeError)    // invalid: doesn't start from "/"
		So(ctx.AddMetric("/plugin/group2/[subgroup2=id12]/metric1", 20), ShouldBeError) // invalid: using dyn. element
	})

	return nil
}

func (s *SuiteT) TestWithoutDefinitionCollector() {
	logrus.SetLevel(logrus.TraceLevel)

	// Arrange
	jsonConfig := []byte(`{}`)
	mtsSelector := []string{
		"/plugin/group1/subgroup1/metric1",
		"/plugin/group2/{id.*}/metric1",
		"/plugin/group3/subgroup3/{.*}",
		"/plugin/group3/subgroup4/**",
	}

	noDefCollector := &noDefinitionCollector{t: s.T()}
	ln := s.startCollector(noDefCollector)
	s.startClient(ln.Addr().String())

	Convey("Validate that collector can gather metric when no definition is provided", s.T(), func() {
		// Act
		_, err := s.sendLoad("task-1", jsonConfig, mtsSelector)
		So(err, ShouldBeNil)

		collMts, err := s.sendCollect("task-1")
		So(err, ShouldBeNil)
		So(len(collMts.MetricSet), ShouldEqual, 5)

		// Validate modifiers
		expectedCustomTimestamp := time.Date(2020, 10, 13, 12, 13, 14, 0, time.UTC)
		So(collMts.MetricSet[0].Description, ShouldEqual, "metric1 custom description")
		So(collMts.MetricSet[1].Unit, ShouldEqual, "MU")
		So(collMts.MetricSet[1].Timestamp.Sec, ShouldEqual, expectedCustomTimestamp.Unix())
		So(collMts.MetricSet[1].Timestamp.Nsec, ShouldEqual, int64(expectedCustomTimestamp.Nanosecond()))
		So(collMts.MetricSet[2].Tags, ShouldResemble, map[string]string{"a": "b", "c": "d"})

		time.Sleep(2 * time.Second)
		_, err = s.sendKill()
		So(err, ShouldBeNil)
	})
}

/*****************************************************************************/

func (s *SuiteT) TestUnloadingRunningCollector() {
	// Arrange
	jsonConfig := []byte(`{}`)
	var mtsSelector []string

	longRunningCollector := &longRunningCollector{collectDuration: 30 * time.Second}
	ln := s.startCollector(longRunningCollector)
	s.startClient(ln.Addr().String())

	errCollectCh := make(chan error, 1)

	Convey("Validate ability to kill collector in case processing takes too much time", s.T(), func() {

		// Act - collect is processing for 30 seconds, but unload comes right after request. Should unblock after 10s with error.
		go func() {
			_, _ = s.sendLoad("task-1", jsonConfig, mtsSelector)
			_, err := s.sendCollect("task-1")
			errCollectCh <- err
		}()

		Convey("Client is able to send unload request and receive no-error response", func() {
			// Act
			time.Sleep(5 * time.Second) // Delay needed to be sure that sendLoad() and sendCollect() in goroutine above were requested
			unloadResp, unloadErr := s.sendUnload("task-1")

			// Assert (kill response)
			So(unloadErr, ShouldBeNil)
			So(unloadResp, ShouldNotBeNil)

			// Assert (plugin has stopped working)
			select {
			case <-errCollectCh:
				// ok
			case <-time.After(expectedUnloadTimeout):
				s.T().Fatal("plugin should have been ended")
			}

			// Assert that Collect was called
			So(longRunningCollector.collectCalls, ShouldEqual, 1)
		})
	})
}

/*****************************************************************************/

type streamingCollector struct {
	collectCalls int
	completed    bool
}

func (c *streamingCollector) StreamingCollect(ctx plugin.CollectContext) error {
	c.collectCalls++
	for {
		select {
		case <-ctx.Done():
			c.completed = true
			return nil
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (s *SuiteT) TestUnloadingRunningStreaming() {
	// Arrange
	jsonConfig := []byte(`{}`)
	var mtsSelector []string

	streamingCollector := &streamingCollector{}
	ln := s.startStreamingCollector(streamingCollector)
	s.startClient(ln.Addr().String())

	errCollectCh := make(chan error, 1)

	Convey("Validate ability to kill collector in case processing takes too much time", s.T(), func() {

		// Act - streaming is processing, but unload comes right after request. Should unblock after 10s with error.
		go func() {
			_, _ = s.sendLoad("task-1", jsonConfig, mtsSelector)
			_, err := s.sendCollect("task-1")
			errCollectCh <- err
		}()

		Convey("Client is able to send unload request and receive no-error response", func() {
			// Act
			time.Sleep(5 * time.Second) // Delay needed to be sure that sendLoad() and sendCollect() in goroutine above were requested
			unloadResp, unloadErr := s.sendUnload("task-1")

			// Assert (kill response)
			So(unloadErr, ShouldBeNil)
			So(unloadResp, ShouldNotBeNil)

			// Assert (plugin has stopped working)
			select {
			case <-errCollectCh:
				// ok
			case <-time.After(expectedUnloadTimeout):
				s.T().Fatal("plugin should have been ended")
			}

			// Assert that Collect was called
			So(streamingCollector.collectCalls, ShouldEqual, 1)
		})
	})
}

/*****************************************************************************/

type collectWithAlwaysApply struct {
	t *testing.T
}

func (c *collectWithAlwaysApply) Collect(ctx plugin.CollectContext) error {
	Convey("Validate AlwaysApply return values", c.t, func() {
		sat, err1 := ctx.AlwaysApply("/coll/group1/*", plugin.MetricTag("ka", "va"))
		So(err1, ShouldBeNil)

		// Should apply tag ka: va
		ctx.AddMetric("/coll/group1/metric1", 11, plugin.MetricTag("k1", "v1")) // mts.MetricSet[0]
		ctx.AddMetric("/coll/group1/metric2", 12, plugin.MetricTag("k2", "v2")) // mts.MetricSet[1]
		ctx.AddMetric("/coll/group2/metric3", 13, plugin.MetricTag("k3", "v3")) // mts.MetricSet[2]

		sat.Dismiss()

		// Should not more apply tag ka: va
		ctx.AddMetric("/coll/group1/metric1", 21, plugin.MetricTag("k1", "v1")) // mts.MetricSet[3]
		ctx.AddMetric("/coll/group1/metric2", 22, plugin.MetricTag("k2", "v2")) // mts.MetricSet[4]
		ctx.AddMetric("/coll/group2/metric3", 23, plugin.MetricTag("k3", "v3")) // mts.MetricSet[5]

		sat2, err2 := ctx.AlwaysApply("/coll/group3/metric4", plugin.MetricTag("kb", "vb"))
		sat3, err3 := ctx.AlwaysApply("/coll/group3/*", plugin.MetricTag("kc", "vc"))
		So(err2, ShouldBeNil)
		So(err3, ShouldBeNil)

		// Should apply tag kb: vb and kc: vc
		ctx.AddMetric("/coll/group3/metric4", 31) // mts.MetricSet[6]

		// Should apply kc: vc
		ctx.AddMetric("/coll/group3/metric5", 41) // mts.MetricSet[7]

		sat3.Dismiss()

		// Should apply tag kb: vb
		ctx.AddMetric("/coll/group3/metric4", 51) // mts.MetricSet[8]

		sat2.Dismiss()

		// Shouldn't apply any tag
		ctx.AddMetric("/coll/group3/metric4", 61) // mts.MetricSet[9]

		// This one shouldn't apply in the next collect
		_, err4 := ctx.AlwaysApply("/coll/**", plugin.MetricTag("kg", "vg"))
		So(err4, ShouldBeNil)
	})

	return nil
}

func (s *SuiteT) TestCollectorWithAlwaysApply() {
	// Arrange
	const collectNumber = 2 // test two consecutive collect

	jsonConfig := []byte(`{}`)
	mtsSelector := []string{}

	collector := &collectWithAlwaysApply{t: s.T()}
	ln := s.startCollector(collector)
	s.startClient(ln.Addr().String())

	Convey("Validate collector can utilize method AlwaysApply", s.T(), func() {
		_, _ = s.sendLoad("task-1", jsonConfig, mtsSelector)

		for i := 0; i < collectNumber; i++ {
			Convey(fmt.Sprintf("Collect no. %d", i+1), func() {
				mts, err := s.sendCollect("task-1")

				So(err, ShouldBeNil)
				So(mts.MetricSet, ShouldNotBeNil)
				So(len(mts.MetricSet), ShouldEqual, 10)

				So(mts.MetricSet[0].Tags, ShouldResemble, map[string]string{"k1": "v1", "ka": "va"})
				So(mts.MetricSet[1].Tags, ShouldResemble, map[string]string{"k2": "v2", "ka": "va"})
				So(mts.MetricSet[2].Tags, ShouldResemble, map[string]string{"k3": "v3"})
				So(mts.MetricSet[3].Tags, ShouldResemble, map[string]string{"k1": "v1"})
				So(mts.MetricSet[4].Tags, ShouldResemble, map[string]string{"k2": "v2"})
				So(mts.MetricSet[5].Tags, ShouldResemble, map[string]string{"k3": "v3"})

				So(mts.MetricSet[6].Tags, ShouldResemble, map[string]string{"kb": "vb", "kc": "vc"})
				So(mts.MetricSet[7].Tags, ShouldResemble, map[string]string{"kc": "vc"})
				So(mts.MetricSet[8].Tags, ShouldResemble, map[string]string{"kb": "vb"})
				So(mts.MetricSet[9].Tags, ShouldBeNil)
			})
		}
	})
}

/*****************************************************************************/

type collectorWithOTELMetrics struct {
	t *testing.T
}

func (c *collectorWithOTELMetrics) Collect(ctx plugin.CollectContext) error {
	var err error

	Convey("Validate AddMetrics won't return errors", c.t, func() {
		err = ctx.AddMetric("/coll/otel/default", 10)
		So(err, ShouldBeNil)

		err = ctx.AddMetric("/coll/otel/gauge", 10.5, plugin.MetricTypeGauge())
		So(err, ShouldBeNil)

		err = ctx.AddMetric("/coll/otel/counter", 67, plugin.MetricTypeCounter())
		So(err, ShouldBeNil)

		counter := plugin.Summary{
			Count: 14,
			Sum:   3.54,
		}

		err = ctx.AddMetric("/coll/otel/summary", counter, plugin.MetricTypeSummary())
		So(err, ShouldBeNil)

		err = ctx.AddMetric("/coll/otel/summary_ptr", &counter, plugin.MetricTypeSummary())
		So(err, ShouldBeNil)

		histogram := plugin.Histogram{
			DataPoints: map[float64]float64{
				0.10:        10,
				0.20:        20,
				0.50:        25,
				1:           10,
				5:           25,
				10:          50,
				math.Inf(1): 100,
			},
			Count: 10,
			Sum:   50,
		}

		err = ctx.AddMetric("/coll/otel/histogram", histogram, plugin.MetricTypeHistogram())
		So(err, ShouldBeNil)

		err = ctx.AddMetric("/coll/otel/histogram_ptr", &histogram, plugin.MetricTypeHistogram())
		So(err, ShouldBeNil)
	})

	return nil
}

func (s *SuiteT) TestCollectingOTELTypes() {
	// Arrange
	jsonConfig := []byte(`{}`)
	var mtsSelector []string

	collector := &collectorWithOTELMetrics{t: s.T()}
	ln := s.startCollector(collector)
	s.startClient(ln.Addr().String())

	Convey("Validate collector can collect OTEL-types metrics", s.T(), func() {
		_, _ = s.sendLoad("task-1", jsonConfig, mtsSelector)

		mts, err := s.sendCollect("task-1")

		So(err, ShouldBeNil)
		So(mts.MetricSet, ShouldNotBeNil)
		So(len(mts.MetricSet), ShouldEqual, 7)

		So(mts.MetricSet[0].Type, ShouldEqual, pluginrpc.MetricType_UNKNOWN)
		So(mts.MetricSet[0].Value.GetVInt64(), ShouldEqual, 10)

		So(mts.MetricSet[1].Type, ShouldEqual, pluginrpc.MetricType_GAUGE)
		So(mts.MetricSet[1].Value.GetVDouble(), ShouldEqual, 10.5)

		So(mts.MetricSet[2].Type, ShouldEqual, pluginrpc.MetricType_COUNTER)
		So(mts.MetricSet[2].Value.GetVInt64(), ShouldEqual, 67)

		So(mts.MetricSet[3].Type, ShouldEqual, pluginrpc.MetricType_SUMMARY)
		So(mts.MetricSet[3].Value.GetVSummary().Sum, ShouldEqual, 3.54)
		So(mts.MetricSet[3].Value.GetVSummary().Count, ShouldEqual, 14)

		So(mts.MetricSet[4].Type, ShouldEqual, pluginrpc.MetricType_SUMMARY)
		So(mts.MetricSet[4].Value.GetVSummary().Sum, ShouldEqual, 3.54)
		So(mts.MetricSet[4].Value.GetVSummary().Count, ShouldEqual, 14)

		So(mts.MetricSet[5].Type, ShouldEqual, pluginrpc.MetricType_HISTOGRAM)
		So(mts.MetricSet[5].Value.GetVHistogram().Sum, ShouldEqual, 50)
		So(mts.MetricSet[5].Value.GetVHistogram().Count, ShouldEqual, 10)
		So(mts.MetricSet[5].Value.GetVHistogram().Bounds, ShouldResemble, []float64{0.10, 0.20, 0.50, 1, 5, 10, math.Inf(1)})
		So(mts.MetricSet[5].Value.GetVHistogram().Values, ShouldResemble, []float64{10, 20, 25, 10, 25, 50, 100})

		So(mts.MetricSet[6].Type, ShouldEqual, pluginrpc.MetricType_HISTOGRAM)
		So(mts.MetricSet[6].Value.GetVHistogram().Sum, ShouldEqual, 50)
		So(mts.MetricSet[6].Value.GetVHistogram().Count, ShouldEqual, 10)
		So(mts.MetricSet[6].Value.GetVHistogram().Bounds, ShouldResemble, []float64{0.10, 0.20, 0.50, 1, 5, 10, math.Inf(1)})
		So(mts.MetricSet[6].Value.GetVHistogram().Values, ShouldResemble, []float64{10, 20, 25, 10, 25, 50, 100})
	})
}
