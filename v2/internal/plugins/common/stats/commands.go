package stats

import (
	"time"
)

///////////////////////////////////////////////////////////////////////////////

type StatCommand interface {
	ApplyStat()
}

///////////////////////////////////////////////////////////////////////////////

type loadTaskStat struct {
	sm      *StatisticsController
	taskID  string
	config  string
	filters []string
}

func (ts *loadTaskStat) ApplyStat() {
	ts.sm.applyLoadStat(ts.taskID, ts.config, ts.filters)
}

///////////////////////////////////////////////////////////////////////////////

type unloadTaskStat struct {
	sm     *StatisticsController
	taskID string
}

func (ts *unloadTaskStat) ApplyStat() {
	ts.sm.applyUnloadStat(ts.taskID)
}

///////////////////////////////////////////////////////////////////////////////

type collectTaskStat struct {
	sm           *StatisticsController
	taskID       string
	metricsCount int
	success      bool
	startTime    time.Time
	processTime  time.Time
}

func (ts *collectTaskStat) ApplyStat() {
	ts.sm.applyCollectStat(ts.taskID, ts.metricsCount, ts.success, ts.startTime, ts.processTime)
}

///////////////////////////////////////////////////////////////////////////////

type streamTaskStat struct {
	sm           *StatisticsController
	taskID       string
	metricsCount int
	startTime    time.Time
	lastUpdate   time.Time
}

func (ts *streamTaskStat) ApplyStat() {
	ts.sm.applyStreamStat(ts.taskID, ts.metricsCount, ts.startTime, ts.lastUpdate)
}
