package stats

import (
	"encoding/json"
	"time"
)

type Statistics struct {
	PluginInfo   pluginInfoFields
	Tasks        tasksFields
	TasksDetails map[int]taskDetailsFields
}

/*****************************************************************************/

type pluginInfoFields struct {
	Name      string
	Version   string
	StartTime time.Time
	Options   string // todo: format?
}

type tasksFields struct {
	CurrentlyActiveTasks int
	TotalActiveTasks     int
	TotalCollectsRequest int

	AvgProcessingTime time.Duration
	MaxProcessingTime time.Duration

	totalProcessingTime time.Duration
}

type taskDetailsFields struct {
	Configuration json.RawMessage
	Filters       []string

	LoadedTime time.Time

	CollectRequest       int
	TotalMetrics         int
	AvgMetricsPerCollect int

	TotalProcessingTime string
	AvgProcessingTime   string
	MaxProcessingTime   string

	totalProcessingTime time.Duration
	avgProcessingTime   time.Duration
	maxProcessingTime   time.Duration
}
