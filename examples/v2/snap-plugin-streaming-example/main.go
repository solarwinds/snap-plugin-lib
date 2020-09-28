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
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
	"github.com/solarwinds/snap-plugin-lib/v2/runner"
)

const (
	pluginName    = "example-streaming"
	pluginVersion = "0.0.1"

	maxProbeDuration = 1 * time.Second
)

var log = logrus.WithFields(logrus.Fields{
	"layer": "plugin",
	"name":  pluginName,
})

type streamCollector struct {
	probeID int
}

func (c *streamCollector) StreamingCollect(ctx plugin.CollectContext) error {
	for {
		select {
		case <-ctx.Done():
			log.Info("Handling end of stream")
			return nil
		case <-time.After(time.Duration(rand.Intn(int(maxProbeDuration)))):
			log.Debug("Gathering metric")

			c.probeID++
			_ = ctx.AddMetric("/stream/probes/id", c.probeID)
		}
	}
}

func main() {
	runner.StartStreamingCollectorWithContext(context.Background(), &streamCollector{}, pluginName, pluginVersion)
}
