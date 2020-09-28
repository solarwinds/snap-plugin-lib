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

func (s simpleCollector) Collect(ctx plugin.CollectContext) error {
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
	runner.StartCollectorWithContext(context.Background(), &simpleCollector{}, "example", "1.0.0")
}
