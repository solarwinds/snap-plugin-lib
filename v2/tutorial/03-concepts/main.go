/*
 Copyright (c) 2020 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package main

import (
	"context"
	"time"

	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
	"github.com/solarwinds/snap-plugin-lib/v2/runner"
)

type simpleCollector struct{}

func (s simpleCollector) format(ctx plugin.Context) string {
	fm, _ := ctx.ConfigValue("format")
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

	_ = ctx.AddMetric("/example/date/day", t.Day(),
		plugin.MetricTimestamp(time.Now().Add(2*time.Hour)))

	_ = ctx.AddMetric("/example/time/hour", hour,
		plugin.MetricDescription("custom description for an hour metric"),
		plugin.MetricUnit("HH"))

	// Count metrics
	startTime, _ := ctx.Load("startTime")
	runningDuration := int(time.Since(startTime.(time.Time)).Seconds())
	_ = ctx.AddMetric("/example/count/running", runningDuration)

	return nil
}

func main() {
	runner.StartCollectorWithContext(context.Background(), &simpleCollector{}, "example", "1.0.0")
}
