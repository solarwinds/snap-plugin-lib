package stats

import (
	"time"
)

type Statistics struct {
	pluginInfo   pluginInfoFields
	tasks        tasksFields
	tasksDetails map[int]taskDetailsFields
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
	Configuration string
	Filters       []string

	LoadedTime           time.Time
	CollectRequest       int
	TotalMetrics         int
	AvgMetricsPerCollect int

	TotalProcessingTime time.Duration
	AvgProcessingTime   time.Duration
	MaxProcessingTime   time.Duration
}
