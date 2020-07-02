package main

import (
	"context"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
)

type simpleCollector struct{}

func (s simpleCollector) PluginDefinition(def plugin.CollectorDefinition) error {
	cfg := `
format: short # format of hour (short 0-12, long 0-24)
options:
  - zone: UTC # time zone
`
	_ = def.DefineExampleConfig(cfg)

	def.DefineMetric("/example/date/day", "", true, "Current day")
	def.DefineMetric("/example/date/month", "", true, "Current month")
	def.DefineMetric("/example/time/hour", "h", true, "Current hour")
	def.DefineMetric("/example/time/minute", "m", true, "Current minute")
	def.DefineMetric("/example/time/second", "s", true, "Current second")
	def.DefineMetric("/example/count/running", "s", false, "Time since task was loaded")

	return nil
}

func (s simpleCollector) format(ctx plugin.Context) string {
	fm, _ := ctx.Config("format")
	if fm == "short" {
		return fm
	}
	return "long"
}

func (s simpleCollector) Load(ctx plugin.Context) error {
	ctx.Store("startTime", time.Now())
	return nil
}

func (s simpleCollector) Collect(ctx plugin.CollectContext) error {
	// Collect data
	t := time.Now()

	// Handle configuration
	hour := t.Hour()
	if s.format(ctx) == "short" {
		hour %= 12
	}

	// Convert data to metric form
	_ = ctx.AddMetric("/example/date/day", t.Day(), plugin.MetricTag("weekday", t.Weekday().String()))
	_ = ctx.AddMetric("/example/date/month", int(t.Month()))
	_ = ctx.AddMetric("/example/time/hour", hour)
	_ = ctx.AddMetric("/example/time/minute", t.Minute())
	_ = ctx.AddMetric("/example/time/second", t.Second())

	// Count metrics
	startTime, _ := ctx.Load("startTime")
	runningDuration := int(time.Since(startTime.(time.Time)).Seconds())
	_ = ctx.AddMetric("/example/count/running", runningDuration)

	return nil
}

func main() {
	runner.StartCollectorWithContext(context.Background(), &simpleCollector{}, "example", "1.0.0")
}
