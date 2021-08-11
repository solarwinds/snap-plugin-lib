/*
 Copyright (c) 2021 SolarWinds Worldwide, LLC

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
	"fmt"

	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
	"github.com/solarwinds/snap-plugin-lib/v2/runner"
)

const (
	pluginName    = "example"
	pluginVersion = "0.0.1"
)

type myCollector struct {
}

func (c *myCollector) Collect(ctx plugin.CollectContext) error {
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			ns := fmt.Sprintf("/example/group_%d/metric_%d", i, j)
			tags := map[string]string{
				"t1": fmt.Sprintf("%d", i),
				"t2": fmt.Sprintf("%d", j),
			}
			_ = ctx.AddMetric(ns, i*j, plugin.MetricTags(tags))
		}
	}

	return nil
}

func main() {
	runner.StartCollectorWithContext(context.Background(), &myCollector{}, "example", pluginVersion)
}
