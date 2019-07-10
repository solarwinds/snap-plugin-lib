package stats

import (
	"encoding/json"
	"time"
)

type Statistics struct {
	PluginInfo   pluginInfo
	TasksSummary tasksSummary
	TasksDetails map[int]taskDetails
}

/*****************************************************************************/

type pluginInfo struct {
	Name           string
	Version        string
	StartTime      time.Time
	OperatingTime  string
	CmdLineOptions string
	Options        json.RawMessage
}

type tasksSummary struct {
	Counters        summaryCounters
	ProcessingTimes processingTimes
}

type taskDetails struct {
	Configuration json.RawMessage
	Filters       []string

	Counters        tasksCounters
	Loaded          eventTimes
	ProcessingTimes processingTimes
}

///////////////////////////////////////////////////////////////////////////////

type summaryCounters struct {
	CurrentlyActiveTasks int
	TotalActiveTasks     int
	TotalCollectsRequest int
}

type tasksCounters struct {
	CollectRequests      int
	TotalMetrics         int
	AvgMetricsPerCollect int
}

///////////////////////////////////////////////////////////////////////////////

type processingTimes struct {
	Total   time.Duration
	Average time.Duration
	Maximum time.Duration
}

type processingTimesJSON struct {
	Total   string
	Average string
	Maximum string
}

func (pt processingTimes) MarshalJSON() ([]byte, error) {
	ptJSON := processingTimesJSON{
		Total:   pt.Total.String(),
		Average: pt.Average.String(),
		Maximum: pt.Maximum.String(),
	}

	return json.Marshal(ptJSON)
}

///////////////////////////////////////////////////////////////////////////////

type eventTimes struct {
	Time time.Time
	Ago  time.Duration
}

func (ot eventTimes) MarshalJSON() ([]byte, error) {
	otJSON := operatingTimesJSON{
		Time: ot.Time.Format(time.StampMicro),
		Ago:  time.Since(ot.Time).String(),
	}

	return json.Marshal(otJSON)
}

type operatingTimesJSON struct {
	Time string
	Ago  string
}
