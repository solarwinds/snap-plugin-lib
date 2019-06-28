package main

import (
	"math/rand"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{
	"layer": "plugin",
	"name":  "example-collector",
})

type myCollector struct {
}

func (*myCollector) DefineMetrics(plugin.CollectorDefinition) error {
	return nil
}

func (*myCollector) Collect(ctx plugin.Context) error {
	log.Trace("Collect executed")

	ctx.AddMetric("/example/demo/random1", rand.Intn(10))
	ctx.AddMetric("/example/demo/random2", rand.Intn(20))

	return nil
}

func (*myCollector) Load(ctx plugin.Context) error {
	log.Tracef("Load executed, configs=%s", ctx.ConfigKeys())

	return nil
}

func (*myCollector) Unload(ctx plugin.Context) error {
	log.Trace("Unload executed")

	return nil
}

func main() {
	runner.StartCollector(&myCollector{}, "example-collector", "0.0.1")
}
