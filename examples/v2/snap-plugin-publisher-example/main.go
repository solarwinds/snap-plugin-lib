package main

import (
	"github.com/sirupsen/logrus"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
)

const (
	pluginName    = "example"
	pluginVersion = "0.0.1"
)

type myPublisher struct {
}

func (m myPublisher) PluginDefinition(def plugin.PublisherDefinition) error {
	_ = def.DefineTasksPerInstanceLimit(1)
	_ = def.DefineInstancesLimit(4)
	return nil
}

func (m myPublisher) Publish(ctx plugin.PublishContext) error {
	logrus.Infof("Number of metrics: %v\n", ctx.Count())

	for _, mt := range ctx.ListAllMetrics() {
		logrus.Infof(" - %s=%v [%v]\n", mt.Namespace(), mt.Value(), mt.Tags())
	}

	return nil
}

func main() {
	runner.StartPublisher(&myPublisher{}, pluginName, pluginVersion)
}
