package main

import (
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/runner"
)

type simpleCollector struct{}

func (s simpleCollector) Collect(ctx plugin.Context) error {
	// Collect data
	t := time.Now()

	// Convert data to metric form
	_ = ctx.AddMetric("/example/date/day", t.Day())
	_ = ctx.AddMetric("/example/date/month", int(t.Month()))
	_ = ctx.AddMetric("/example/time/hour", t.Hour())
	_ = ctx.AddMetric("/example/time/minute", t.Minute())
	_ = ctx.AddMetric("/example/time/second", t.Second())

	return nil
}

func main() {
	runner.StartCollector(&simpleCollector{}, "example", "1.0.0")
}
