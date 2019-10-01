package collector

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

const (
	defaultMinCpuUsage             = 0.05
	defaultMinMemoryUsage          = 0.01
	defaultTotalCpuMeasureDuration = "1s"

	configObjectKey = "config"
)

type config struct {
	Processes               configProcesses
	TotalCpuMeasureDuration string
}

type configProcesses struct {
	MinCpuUsage    float64
	MinMemoryUsage float64
}

func defaultConfig() *config {
	return &config{
		Processes: configProcesses{
			MinCpuUsage:    defaultMinCpuUsage,
			MinMemoryUsage: defaultMinMemoryUsage,
		},
		TotalCpuMeasureDuration: defaultTotalCpuMeasureDuration,
	}
}

func handleConfig(ctx plugin.Context) error {
	cfg := defaultConfig()

	err := json.Unmarshal(ctx.RawConfig(), cfg)
	if err != nil {
		return fmt.Errorf("invalid config: %v", err)
	}

	_, err = time.ParseDuration(cfg.TotalCpuMeasureDuration)
	if err != nil {
		return fmt.Errorf("invalid value for totalCpuMeasureDuration: %v", err)
	}

	if cfg.Processes.MinCpuUsage < 0 || cfg.Processes.MinCpuUsage > 100 {
		return fmt.Errorf("invalid value for minCpuUsage: %v", err)
	}

	if cfg.Processes.MinMemoryUsage < 0 || cfg.Processes.MinMemoryUsage > 100 {
		return fmt.Errorf("invalid value for minMemoryUsage: %v", err)
	}

	ctx.Store(configObjectKey, cfg)

	return nil
}

func getConfig(ctx plugin.Context) *config {
	obj, ok := ctx.Load(configObjectKey)
	if !ok {
		return defaultConfig()
	}
	return obj.(*config)
}
