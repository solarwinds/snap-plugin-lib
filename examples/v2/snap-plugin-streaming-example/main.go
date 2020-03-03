package main

import (
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

const (
	pluginName    = "example-streaming"
	pluginVersion = "0.0.1"

	maxProbeDuration = 10 * time.Second
)

var log = logrus.WithFields(logrus.Fields{
	"layer": "plugin",
	"name":  pluginName,
})

type streamCollector struct {
	probeID int
}

func (c *streamCollector) StreamingCollect(ctx plugin.CollectContext) {
	log.Trace("StreamingCollect start")

	c.probeID++
	_ = ctx.AddMetric("/stream/probes/id", c.probeID)

	waitDuration := time.Duration(rand.Intn(int(maxProbeDuration)))
	time.Sleep(waitDuration)
}

func main() {
	runner.StartStreamingCollector(&streamCollector{}, pluginName, pluginVersion)
}
