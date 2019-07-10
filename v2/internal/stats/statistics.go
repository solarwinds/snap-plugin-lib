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
	Name           string
	Version        string
	StartTime      time.Time
	OperatingTime  string
	CmdLineOptions string
	Options        json.RawMessage
}

type tasksFields struct {
	CurrentlyActiveTasks int
	TotalActiveTasks     int
	TotalCollectsRequest int

	AvgProcessingTime string
	MaxProcessingTime string

	avgProcessingTime time.Duration
	maxProcessingTime time.Duration

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
