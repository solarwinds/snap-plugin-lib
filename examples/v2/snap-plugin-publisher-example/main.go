package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
	"github.com/sirupsen/logrus"
)

const (
	pluginName    = "example"
	pluginVersion = "0.0.1"
)

type myPublisher struct {
}

func (m myPublisher) PluginDefinition(def plugin.PublisherDefinition) error {
	_ = def.DefineTasksPerInstanceLimit(4)
	_ = def.DefineInstancesLimit(3)
	return nil
}

func (m myPublisher) Publish(ctx plugin.PublishContext) error {
	logrus.Infof("Number of metrics: %v\n", ctx.Count())

	// Simulate publisher processing
	time.Sleep(time.Duration(rand.Intn(int(1 * time.Second))))

	for _, mt := range ctx.ListAllMetrics() {
		logrus.Infof(" - %s=%v [%v]\n", mt.Namespace(), mt.Value(), mt.Tags())
	}

	return nil
}

func main() {
	runner.StartPublisher(context.Background(), &myPublisher{}, pluginName, pluginVersion)
}
