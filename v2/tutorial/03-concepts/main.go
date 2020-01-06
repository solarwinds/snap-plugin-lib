package main

import (
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
)

type simpleCollector struct{}

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
	_ = ctx.AddMetric("/example/date/day", t.Day())
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
	runner.StartCollector(&simpleCollector{}, "example", "1.0.0")
}
