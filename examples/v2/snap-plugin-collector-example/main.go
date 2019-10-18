package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
	"github.com/sirupsen/logrus"
)

const (
	pluginName    = "example"
	pluginVersion = "0.0.1"

	maxCollectDuration = 5 * time.Second
)

var log = logrus.WithFields(logrus.Fields{
	"layer": "plugin",
	"name":  pluginName,
})

type myCollector struct {
}

var exampleConfig = `
# Random value used to calculate metric4
crand: 40

# other tree-like configuration
credentials:
  user: admin
  password: secure1
  token: abcd-1234
`

func (*myCollector) PluginDefinition(def plugin.CollectorDefinition) error {
	_ = def.DefineTasksPerInstanceLimit(5)
	_ = def.DefineInstancesLimit(plugin.NoLimit)

	def.DefineMetric("/example/static/random1", "", true, "Random value (0-10)")
	def.DefineMetric("/example/static/random2", "", true, "Random value (0-20)")
	def.DefineMetric("/example/global/random3", "", true, "Random value (0-50)")
	def.DefineMetric("/example/config/random4", "", true, "Random value (0-crand)")

	def.DefineMetric("/example/nodefault/random5", "", false, "Random value")

	_ = def.DefineExampleConfig(exampleConfig)

	return nil
}

func (*myCollector) Collect(ctx plugin.CollectContext) error {
	log.Trace("Collect executed")

	// Simulate collector processing
	time.Sleep(time.Duration(rand.Intn(int(maxCollectDuration))))

	_ = ctx.AddMetric("/example/static/random1", rand.Intn(10))
	_ = ctx.AddMetric("/example/static/random2", rand.Intn(20))

	globalRandom, _ := ctx.Load("random")
	_ = ctx.AddMetric("/example/global/random3", globalRandom)

	configRandom, ok := ctx.Load("random-config")
	if ok {
		_ = ctx.AddMetricWithTags("/example/config/random4", rand.Intn(configRandom.(int)),
			map[string]string{"random": fmt.Sprintf("%v", configRandom)})
	}

	return nil
}

func (*myCollector) Load(ctx plugin.Context) error {
	log.Tracef("Load executed, configs=%s", ctx.ConfigKeys())

	ctx.Store("random", rand.Intn(50))

	v, ok := ctx.Config("crand")
	if ok {
		intV, err := strconv.Atoi(v)
		if err == nil {
			ctx.Store("random-config", intV)
		}
	}

	return nil
}

func (*myCollector) Unload(ctx plugin.Context) error {
	log.Trace("Unload executed")

	return nil
}

func main() {
	runner.StartCollector(&myCollector{}, pluginName, pluginVersion)
}
