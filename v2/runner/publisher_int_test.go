//go:build medium
// +build medium

/*
 Copyright (c) 2021 SolarWinds Worldwide, LLC

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
	"net"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	collProxy "github.com/solarwinds/snap-plugin-lib/v2/internal/plugins/collector/proxy"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/plugins/common/stats"
	pubProxy "github.com/solarwinds/snap-plugin-lib/v2/internal/plugins/publisher/proxy"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/service"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/types"
	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
	"github.com/solarwinds/snap-plugin-lib/v2/pluginrpc"
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
		contextManager := collProxy.NewContextManager(context.Background(), types.NewCollector("test-collector", "1.0.0", collector), statsController)
		service.StartCollectorGRPC(context.Background(), grpc.NewServer(), contextManager, ln, 0, 0)
		s.endControllerCh <- true
	}()

	return ln
}

func (s *PublisherMediumSuite) startPublisher(publisher plugin.Publisher) net.Listener {
	var ln net.Listener

	s.startedPublisher = publisher
	ln, _ = net.Listen("tcp", "127.0.0.1:")

	go func() {
		statsController := &stats.EmptyController{}
		contextManager := pubProxy.NewContextManager(publisher, statsController)
		service.StartPublisherGRPC(context.Background(), grpc.NewServer(), contextManager, ln, 0, 0)
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

	var mts []*pluginrpc.Metric

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
	if err != nil {
		return err
	}

	_, err = stream.CloseAndRecv()
	return err
}

///////////////////////////////////////////////////////////////////////////////

func TestPublisherMedium(t *testing.T) {
	suite.Run(t, new(PublisherMediumSuite))
}

///////////////////////////////////////////////////////////////////////////////

type oneMetricCollector struct {
	collectCalls int
}

func (s *oneMetricCollector) Collect(ctx plugin.CollectContext) error {
	s.collectCalls++
	_ = ctx.AddMetric("/example/group1/metric1", 1, plugin.MetricTags(map[string]string{"k1": "v1", "k3": "v3"}))
	return nil
}

type configurablePublisher struct {
	t            *testing.T
	loadCalls    int
	unloadCalls  int
	publishCalls int
}

func (p *configurablePublisher) Load(ctx plugin.Context) error {
	p.loadCalls++

	ctx.Store("stringValue", "value")
	ctx.Store("intValue", 10)

	Convey("Validate that Publisher has access to context features in Load method", p.t, func() {
		cInterval, ok := ctx.ConfigValue("config.interval")
		So(cInterval, ShouldEqual, "10s")
		So(ok, ShouldBeTrue)

		cDuration, ok := ctx.ConfigValue("config.duration")
		So(cDuration, ShouldEqual, "20s")
		So(ok, ShouldBeTrue)

		_, ok = ctx.ConfigValue("config.timeout")
		So(ok, ShouldBeFalse)
	})

	return nil
}

func (p *configurablePublisher) Unload(ctx plugin.Context) error {
	p.unloadCalls++

	return nil
}

func (p *configurablePublisher) Publish(ctx plugin.PublishContext) error {
	p.publishCalls++

	Convey("Validate that Publisher has access to context features in Collect method", p.t, func() {
		Convey("Metrics API", func() {
			// Act
			mts := ctx.ListAllMetrics()

			// Assert
			So(ctx.Count(), ShouldEqual, 1)
			So(len(mts), ShouldEqual, 1)

			So(mts[0].Tags()["k1"], ShouldEqual, "v1")
			So(mts[0].Tags()["k2"], ShouldNotEqual, "v2")
			So(mts[0].Tags()["k3"], ShouldEqual, "v3")

			So(mts[0].Namespace().HasElementOn("example", 0), ShouldBeTrue)
			So(mts[0].Namespace().HasElementOn("group1", 1), ShouldBeTrue)
			So(mts[0].Namespace().HasElementOn("metric1", 2), ShouldBeTrue)
			So(mts[0].Namespace().HasElementOn("metric2", 2), ShouldBeFalse)
			So(mts[0].Namespace().HasElementOn("metric1", 3), ShouldBeFalse)

			So(mts[0].Namespace().HasElement("example"), ShouldBeTrue)
			So(mts[0].Namespace().HasElement("group1"), ShouldBeTrue)
			So(mts[0].Namespace().HasElement("metric1"), ShouldBeTrue)
			So(mts[0].Namespace().HasElement("metric2"), ShouldBeFalse)
		})

		Convey("Config API", func() {
			cInterval, ok := ctx.ConfigValue("config.interval")
			So(cInterval, ShouldEqual, "10s")
			So(ok, ShouldBeTrue)

			cDuration, ok := ctx.ConfigValue("config.duration")
			So(cDuration, ShouldEqual, "20s")
			So(ok, ShouldBeTrue)

			_, ok = ctx.ConfigValue("config.timeout")
			So(ok, ShouldBeFalse)
		})

		Convey("State API", func() {
			v1, ok := ctx.Load("stringValue")
			So(v1, ShouldEqual, "value")
			So(ok, ShouldBeTrue)

			v2, ok := ctx.Load("intValue")
			So(v2, ShouldEqual, 10)
			So(ok, ShouldBeTrue)
		})
	})

	return nil
}

func (s *PublisherMediumSuite) TestConfigurablePublisher() {
	// Arrange
	publisherConfig := []byte(`{
"config": {
	"interval": "10s",
	"duration": "20s"
}}`)

	collector := &oneMetricCollector{}
	publisher := &configurablePublisher{t: s.T()}

	lnColl := s.startCollector(collector) // collector server (plugin)
	lnPub := s.startPublisher(publisher)  // publisher server (plugin)

	s.startCollectorClient(lnColl.Addr().String()) // collector client (snap)
	s.startPublisherClient(lnPub.Addr().String())  // collector client (snap)

	Convey("Test that publisher can access context methods", s.T(), func() {

		Convey("Loading tasks for collector and publisher", func() {
			_, err := s.sendCollectorLoad("task-collector-1", []byte("{}"), []string{})
			So(err, ShouldBeNil)

			_, err = s.sendPublisherLoad("task-publisher-1", publisherConfig)
			So(err, ShouldBeNil)

			So(collector.collectCalls, ShouldEqual, 0)
			So(publisher.publishCalls, ShouldEqual, 0)
			So(publisher.loadCalls, ShouldEqual, 1)
			So(publisher.unloadCalls, ShouldEqual, 0)

			// validation is also done in Load method of publisher
		})

		Convey("Publisher can process all metrics gathered by collector", func() {
			err := s.requestCollectPublishCycle("task-collector-1", "task-publisher-1")
			So(err, ShouldBeNil)

			So(collector.collectCalls, ShouldEqual, 1)
			So(publisher.publishCalls, ShouldEqual, 1)
			So(publisher.loadCalls, ShouldEqual, 1)
			So(publisher.unloadCalls, ShouldEqual, 0)

			// validation is done in Publish method of publisher
		})

		Convey("Unloading tasks from collector and publisher", func() {
			_, err := s.sendCollectorUnload("task-collector-1")
			So(err, ShouldBeNil)

			_, err = s.sendPublisherUnload("task-publisher-1")
			So(err, ShouldBeNil)

			So(collector.collectCalls, ShouldEqual, 1)
			So(publisher.publishCalls, ShouldEqual, 1)
			So(publisher.loadCalls, ShouldEqual, 1)
			So(publisher.unloadCalls, ShouldEqual, 1)
		})

		Convey("Sending kills requests for collector and publisher", func() {
			err := s.sendKills()
			So(err, ShouldBeNil)
		})

		Convey("Validate that both plugins quit", func() {
			completeCh := make(chan bool, 1)

			go func() {
			loop:
				for i := 0; i < 2; i++ {
					select {
					case <-s.endControllerCh:
						// ok
					case <-s.endPublisherCh:
						// ok
					case <-time.After(3 * time.Second):
						break loop
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

///////////////////////////////////////////////////////////////////////////////

type noConfigCollector struct {
	collectCalls int
}

func (c *noConfigCollector) PluginDefinition(def plugin.CollectorDefinition) error {
	def.DefineMetric("/example/group1/metric1", "b", true, "description 11")
	def.DefineMetric("/example/group1/metric2", "b", true, "description 12")
	def.DefineMetric("/example/group1/metric3", "b", true, "description 13")
	def.DefineMetric("/example/group1/metric4", "b", true, "description 14")
	def.DefineMetric("/example/group1/metric5", "b", true, "description 15")
	def.DefineMetric("/example/group2/metric1", "B", true, "description 21")
	def.DefineMetric("/example/group2/metric2", "B", true, "description 22")

	def.DefineGroup("dyngroup", "dyngroup description")
	def.DefineMetric("/example/group3/[dyngroup]/metric1", "B", true, "description 31")

	return nil
}

func (c *noConfigCollector) Collect(ctx plugin.CollectContext) error {
	c.collectCalls++

	_ = ctx.AddMetric("/example/group1/metric1", 11, plugin.MetricTags(map[string]string{"k1": "v1"}))
	_ = ctx.AddMetric("/example/group1/metric2", 12)
	_ = ctx.AddMetric("/example/group1/metric3", 13)
	_ = ctx.AddMetric("/example/group1/metric4", 14, plugin.MetricTags(map[string]string{"k2": "v2", "k3": "v3"}))
	_ = ctx.AddMetric("/example/group1/metric5", 15)

	_ = ctx.AddMetric("/example/group2/metric1", 21)
	_ = ctx.AddMetric("/example/group2/metric2", 22)

	_ = ctx.AddMetric("/example/group3/[dyngroup=abc]/metric1", 31)

	return nil
}

type simplePublisher struct {
	t            *testing.T
	publishCalls int
}

func (p *simplePublisher) Publish(ctx plugin.PublishContext) error {
	p.publishCalls++

	Convey("Validate that all collected metrics are accessed from Publish", p.t, func() {
		// Act
		mts := ctx.ListAllMetrics()

		// Assert
		So(ctx.Count(), ShouldEqual, 8)
		So(len(mts), ShouldEqual, 8)

		So(mts[0].Namespace().String(), ShouldEqual, "/example/group1/metric1")
		So(mts[0].Value(), ShouldEqual, 11)
		So(mts[0].Tags(), ShouldResemble, map[string]string{"k1": "v1"})
		So(mts[0].Tags()["k1"], ShouldEqual, "v1")
		So(mts[0].Description(), ShouldEqual, "description 11")
		So(mts[0].Unit(), ShouldEqual, "b")

		So(mts[1].Namespace().String(), ShouldEqual, "/example/group1/metric2")
		So(mts[1].Value(), ShouldEqual, 12)
		So(mts[1].Tags(), ShouldResemble, map[string]string{})
		So(mts[1].Description(), ShouldEqual, "description 12")
		So(mts[1].Unit(), ShouldEqual, "b")

		So(mts[2].Namespace().String(), ShouldEqual, "/example/group1/metric3")
		So(mts[2].Value(), ShouldEqual, 13)
		So(mts[2].Tags(), ShouldResemble, map[string]string{})
		So(mts[2].Description(), ShouldEqual, "description 13")
		So(mts[2].Unit(), ShouldEqual, "b")

		So(mts[3].Namespace().String(), ShouldEqual, "/example/group1/metric4")
		So(mts[3].Value(), ShouldEqual, 14)
		So(mts[3].Tags(), ShouldResemble, map[string]string{"k3": "v3", "k2": "v2"})
		So(mts[3].Tags()["k2"], ShouldEqual, "v2")
		So(mts[3].Description(), ShouldEqual, "description 14")
		So(mts[3].Unit(), ShouldEqual, "b")

		So(mts[4].Namespace().String(), ShouldEqual, "/example/group1/metric5")
		So(mts[4].Value(), ShouldEqual, 15)
		So(mts[4].Tags(), ShouldResemble, map[string]string{})
		So(mts[4].Description(), ShouldEqual, "description 15")
		So(mts[4].Unit(), ShouldEqual, "b")

		So(mts[5].Namespace().String(), ShouldEqual, "/example/group2/metric1")
		So(mts[5].Value(), ShouldEqual, 21)
		So(mts[5].Tags(), ShouldResemble, map[string]string{})
		So(mts[5].Description(), ShouldEqual, "description 21")
		So(mts[5].Unit(), ShouldEqual, "B")

		So(mts[6].Namespace().String(), ShouldEqual, "/example/group2/metric2")
		So(mts[6].Value(), ShouldEqual, 22)
		So(mts[6].Tags(), ShouldResemble, map[string]string{})
		So(mts[6].Description(), ShouldEqual, "description 22")
		So(mts[6].Unit(), ShouldEqual, "B")

		So(mts[7].Namespace().String(), ShouldEqual, "/example/group3/[dyngroup=abc]/metric1")
		So(mts[7].Value(), ShouldEqual, 31)
		So(mts[7].Tags(), ShouldResemble, map[string]string{})
		So(mts[7].Description(), ShouldEqual, "description 31")
		So(mts[7].Unit(), ShouldEqual, "B")

		mts[7].Namespace()
		So(mts[7].Namespace().At(2).IsDynamic(), ShouldBeTrue)
		So(mts[7].Namespace().At(2).Name(), ShouldEqual, "dyngroup")
		So(mts[7].Namespace().At(2).Value(), ShouldEqual, "abc")
		So(mts[7].Namespace().At(2).Description(), ShouldEqual, "dyngroup description")
	})

	return nil
}

func (s *PublisherMediumSuite) TestSimplePublisher() {
	// Arrange
	collector := &noConfigCollector{}
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

			So(collector.collectCalls, ShouldEqual, 1)
			So(publisher.publishCalls, ShouldEqual, 1)

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
			loop:
				for i := 0; i < 2; i++ {
					select {
					case <-s.endControllerCh:
						// ok
					case <-s.endPublisherCh:
						// ok
					case <-time.After(3 * time.Second):
						break loop
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
