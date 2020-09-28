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
	"fmt"
	"regexp"
	"strings"

	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
	"github.com/solarwinds/snap-plugin-lib/v2/tutorial/08-collector/collector/proxy"
)

var sanitizeRegex = regexp.MustCompile(`[()[\]{}<> ,./?;':"|^!\-_+\\]+`)

type systemCollector struct {
	proxyCollector proxy.Proxy
}

func New(proxy proxy.Proxy) plugin.Collector {
	return systemCollector{
		proxyCollector: proxy,
	}
}

func (s systemCollector) PluginDefinition(def plugin.CollectorDefinition) error {
	def.DefineGroup("processName", "process name")

	def.DefineMetric("/minisystem/processes/[processName]/cpu", "%", true, "CPU Utilization by current process")
	def.DefineMetric("/minisystem/processes/[processName]/memory", "%", true, "Memory Utilization by current process")
	def.DefineMetric("/minisystem/usage/cpu", "%", true, "Total CPU Utilization")
	def.DefineMetric("/minisystem/usage/memory", "%", true, "Total memory Utilization")

	return nil
}

func (s systemCollector) Collect(ctx plugin.CollectContext) error {
	err := s.collectTotalCPU(ctx)
	if err != nil {
		return err
	}

	err = s.collectTotalMemory(ctx)
	if err != nil {
		return err
	}

	err = s.collectProcessesInfo(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s systemCollector) collectTotalCPU(ctx plugin.CollectContext) error {
	cpu, err := s.proxyCollector.TotalCpuUsage()
	if err != nil {
		return fmt.Errorf("can't create metric for total cpu utilization: %v", err)
	}

	_ = ctx.AddMetric("/minisystem/usage/cpu", cpu)
	return nil
}

func (s systemCollector) collectTotalMemory(ctx plugin.CollectContext) error {
	memory, err := s.proxyCollector.TotalMemoryUsage()

	if err != nil {
		return fmt.Errorf("can't create metric for total memory utilization: %v", err)
	}

	_ = ctx.AddMetric("/minisystem/usage/memory", memory)
	return nil
}

func (s systemCollector) collectProcessesInfo(ctx plugin.CollectContext) error {
	procsInfo, err := s.proxyCollector.ProcessesInfo()
	if err != nil {
		return fmt.Errorf("can't create metrics associated with processes")
	}

	for _, p := range procsInfo {
		pName := s.sanitizeName(p.ProcessName)

		cpuMetricNs := fmt.Sprintf("/minisystem/processes/[processName=%s]/cpu", pName)
		_ = ctx.AddMetric(cpuMetricNs, p.CpuUsage, plugin.MetricTag("PID", fmt.Sprintf("%d", p.PID)))

		memMetricNs := fmt.Sprintf("/minisystem/processes/[processName=%s]/memory", pName)
		_ = ctx.AddMetric(memMetricNs, p.MemoryUsage, plugin.MetricTag("PID", fmt.Sprintf("%d", p.PID)))
	}

	return nil
}

func (s systemCollector) sanitizeName(n string) string {
	return strings.ToLower(sanitizeRegex.ReplaceAllString(n, "_"))
}
