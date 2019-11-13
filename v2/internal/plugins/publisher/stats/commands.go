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
	taskId  string
	config  string
}

func (ts *loadTaskStat) ApplyStat() {
	ts.sm.applyLoadStat(ts.taskId, ts.config)
}

///////////////////////////////////////////////////////////////////////////////

type unloadTaskStat struct {
	sm     *StatisticsController
	taskId string
}

func (ts *unloadTaskStat) ApplyStat() {
	ts.sm.applyUnloadStat(ts.taskId)
}

///////////////////////////////////////////////////////////////////////////////

type collectTaskStat struct {
	sm           *StatisticsController
	taskId       string
	metricsCount int
	success      bool
	startTime    time.Time
	processTime  time.Time
}

func (ts *collectTaskStat) ApplyStat() {
	ts.sm.applyPublishStat(ts.taskId, ts.metricsCount, ts.success, ts.startTime, ts.processTime)
}
