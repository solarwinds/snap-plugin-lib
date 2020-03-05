package main

import (
	"math/rand"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
)

const (
	pluginName    = "example-streaming"
	pluginVersion = "0.0.1"

	maxProbeDuration = 1 * time.Second
)

type streamCollector struct {
	probeID int
}

func (c *streamCollector) StreamingCollect(ctx plugin.CollectContext) {
	c.probeID++
	_ = ctx.AddMetric("/stream/probes/id", c.probeID)

	waitDuration := time.Duration(rand.Intn(int(maxProbeDuration)))
	time.Sleep(waitDuration)
}

func main() {
	runner.StartStreamingCollector(&streamCollector{}, pluginName, pluginVersion)
}
