package collector

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	proxy2 "github.com/librato/snap-plugin-lib-go/v2/tutorial/09-config/collector/proxy"
)

var sanitizeRegex = regexp.MustCompile(`[()[\]{}<> ,./?;':"|^!\-_+\\]+`)

type systemCollector struct {
	proxyCollector proxy2.Proxy
}

func New(proxy proxy2.Proxy) plugin.Collector {
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

func (s systemCollector) Load(ctx plugin.Context) error {
	return handleConfig(ctx)
}

func (s systemCollector) Unload(ctx plugin.Context) error {
	return nil
}

func (s systemCollector) collectTotalCPU(ctx plugin.CollectContext) error {
	cfg := getConfig(ctx)
	measurementDur, _ := time.ParseDuration(cfg.TotalCpuMeasureDuration)

	cpu, err := s.proxyCollector.TotalCpuUsage(measurementDur)
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

	cfg := getConfig(ctx)

	for _, p := range procsInfo {
		pName := s.sanitizeName(p.ProcessName)

		if p.CpuUsage >= cfg.Processes.MinCpuUsage {
			cpuMetricNs := fmt.Sprintf("/minisystem/processes/[processName=%s]/cpu", pName)
			_ = ctx.AddMetricWithTags(cpuMetricNs, p.CpuUsage, map[string]string{"PID": fmt.Sprintf("%d", p.PID)})
		}

		if p.MemoryUsage >= cfg.Processes.MinMemoryUsage {
			memMetricNs := fmt.Sprintf("/minisystem/processes/[processName=%s]/memory", pName)
			_ = ctx.AddMetricWithTags(memMetricNs, p.MemoryUsage, map[string]string{"PID": fmt.Sprintf("%d", p.PID)})
		}
	}

	return nil
}

func (s systemCollector) sanitizeName(n string) string {
	return strings.ToLower(sanitizeRegex.ReplaceAllString(n, "_"))
}
