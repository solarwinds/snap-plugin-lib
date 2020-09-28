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

package collector

import (
	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
	"github.com/solarwinds/snap-plugin-lib/v2/tutorial/07-proxy/collector/proxy"
)

type systemCollector struct {
	proxyCollector proxy.Proxy
}

func (s systemCollector) PluginDefinition(def plugin.CollectorDefinition) error {
	return nil
}

func New(proxy proxy.Proxy) systemCollector {
	return systemCollector{
		proxyCollector: proxy,
	}
}

func (s systemCollector) Collect(plugin.CollectContext) error {
	return nil
}
