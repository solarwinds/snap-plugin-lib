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
	pluginName    = "example"
	pluginVersion = "0.0.1"
)

type myPublisher struct {
}

func (m myPublisher) PluginDefinition(def plugin.PublisherDefinition) error {
	_ = def.DefineTasksPerInstanceLimit(4)
	_ = def.DefineInstancesLimit(3)
	return nil
}

func (m myPublisher) Publish(ctx plugin.PublishContext) error {
	logrus.Infof("Number of metrics: %v\n", ctx.Count())

	// Simulate publisher processing
	time.Sleep(time.Duration(rand.Intn(int(1 * time.Second)))) // #nosec G404

	for _, mt := range ctx.ListAllMetrics() {
		logrus.Infof(" - %s=%v [%v]\n", mt.Namespace(), mt.Value(), mt.Tags())
	}

	return nil
}

func main() {
	runner.StartPublisherWithContext(context.Background(), &myPublisher{}, pluginName, pluginVersion)
}
