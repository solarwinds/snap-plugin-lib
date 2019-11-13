package stats

import (
	"encoding/json"
	"time"
)

type Statistics struct {
	PluginInfo   pluginInfo             `json:"Plugin info"`
	TasksSummary tasksSummary           `json:"Tasks summary"`
	TasksDetails map[string]taskDetails `json:"Task details"`
}

/*****************************************************************************/

type pluginInfo struct {
	Name           string          `json:"Name"`
	Version        string          `json:"Version"`
	CmdLineOptions string          `json:"Command-line options"`
	Options        json.RawMessage `json:"Options"`
	Started        eventTimes      `json:"Started"`
}

type tasksSummary struct {
	Counters        summaryCounters `json:"Counters"`
	ProcessingTimes processingTimes `json:"Processing times"`
}

type taskDetails struct {
	Configuration json.RawMessage `json:"Configuration"`
	Filters       []string        `json:"Requested metrics (filters)"`

	Counters        tasksCounters   `json:"Counters"`
	Loaded          eventTimes      `json:"Loaded"`
	ProcessingTimes processingTimes `json:"Processing times"`
	LastMeasurement measurementInfo `json:"Last publication"`
}

///////////////////////////////////////////////////////////////////////////////

type summaryCounters struct {
	CurrentlyActiveTasks int `json:"Currently active tasks"`
	TotalActiveTasks     int `json:"Total active tasks"`
	TotalPublishRequests int `json:"Total publish requests"`
}

type tasksCounters struct {
	CollectRequests      int `json:"Publish requests"`
	TotalMetrics         int `json:"Total metrics"`
	AvgMetricsPerPublish int `json:"Average metrics / Publish"`
}

type measurementInfo struct {
	Occurred         eventTimes
	Duration         time.Duration
	PublishedMetrics int
}

type measurementInfoJSON struct {
	Occurred         eventTimes `json:"Occurred"`
	Duration         string     `json:"Duration"`
	PublishedMetrics int        `json:"Published metrics"`
}

func (mi measurementInfo) MarshalJSON() ([]byte, error) {
	miJSON := measurementInfoJSON{
		Occurred:         mi.Occurred,
		Duration:         mi.Duration.String(),
		PublishedMetrics: mi.PublishedMetrics,
	}

	return json.Marshal(miJSON)
}

///////////////////////////////////////////////////////////////////////////////

type processingTimes struct {
	Total   time.Duration
	Average time.Duration
	Maximum time.Duration
}

type processingTimesJSON struct {
	Total   string `json:"Total"`
	Average string `json:"Average"`
	Maximum string `json:"Maximum"`
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

type operatingTimesJSON struct {
	Time string `json:"Time"`
	Ago  string `json:"Ago"`
}

func (ot eventTimes) MarshalJSON() ([]byte, error) {
	otJSON := operatingTimesJSON{
		Time: ot.Time.Format(time.StampMicro),
		Ago:  time.Since(ot.Time).String(),
	}

	return json.Marshal(otJSON)
}
