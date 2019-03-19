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
	"google.golang.org/grpc"
)

/*****************************************************************************/

const EXPECTED_GRACEFUL_SHUTDOWN_TIMEOUT = 3 * time.Second
const EXPECTED_FORCE_SHUTDOWN_TIMEOUT = 15 * time.Second

/*****************************************************************************/

type simpleCollector struct {
	collectCalls int
}

func (sc *simpleCollector) Collect(ctx plugin.Context) error {
	sc.collectCalls++
	return nil
}

/*****************************************************************************/

func TestMedium_SimpleCollector(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)

	Convey("Validate ability to connect to simplest collector", t, func() {
		endCh := make(chan bool, 1)
		simpleCollector := &simpleCollector{}
		// Start simple plugin
		go func() {
			contextManager := proxy.NewContextManager(simpleCollector, "simple_collector", "1.0.0")
			rpc.StartGRPCController(contextManager)
			endCh <- true
		}()

		// Connect to plugin
		cc, _ := grpc.Dial("localhost:56789", grpc.WithInsecure())
		defer func() { _ = cc.Close() }()

		coll := rpc.NewCollectorClient(cc)
		cont := rpc.NewControllerClient(cc)

		// Act
		pingResponse, pingErr := sendPing(cont)
		So(pingErr, ShouldBeNil)
		So(pingResponse, ShouldNotBeNil)

		collectResponse, collectErr := sendCollect(coll)
		So(collectErr, ShouldBeNil)
		So(collectResponse, ShouldNotBeNil)
		So(simpleCollector.collectCalls, ShouldEqual, 1)

		for i := 0; i < 5; i++ {
			sendCollect(coll)
		}
		So(simpleCollector.collectCalls, ShouldEqual, 6)

		killResponse, killErr := sendKill(cont)
		So(killErr, ShouldBeNil)
		So(killResponse, ShouldNotBeNil)

		// Assert plugin has completed
		select {
		case <-endCh:
		// ok
		case <-time.After(EXPECTED_GRACEFUL_SHUTDOWN_TIMEOUT):
			t.Fatal("plugin should have been ended")
		}
	})
}

func sendPing(c rpc.ControllerClient) (*rpc.PingResponse, error) {
	response, err := c.Ping(context.Background(), &rpc.PingRequest{})
	return response, err
}

func sendKill(c rpc.ControllerClient) (*rpc.KillResponse, error) {
	response, err := c.Kill(context.Background(), &rpc.KillRequest{})
	return response, err
}

func sendLoad(c rpc.CollectorClient) (*rpc.LoadResponse, error) {
	response, err := c.Load(context.Background(), &rpc.LoadRequest{
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

func sendUnload(c rpc.CollectorClient) (*rpc.UnloadResponse, error) {
	response, err := c.Unload(context.Background(), &rpc.UnloadRequest{
		TaskId: 1,
	})
	return response, err
}

func sendCollect(c rpc.CollectorClient) (*rpc.CollectResponse, error) {
	response, err := c.Collect(context.Background(), &rpc.CollectRequest{
		TaskId: 1,
	})

	// todo: handle streams
	response.Recv()

	return &rpc.CollectResponse{}, err
}
