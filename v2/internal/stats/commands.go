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
	taskId  int
	config  string
	filters []string
}

func (ts *loadTaskStat) ApplyStat() {
	ts.sm.applyLoadStat(ts.taskId, ts.config, ts.filters)
}

///////////////////////////////////////////////////////////////////////////////

type unloadTaskStat struct {
	sm     *StatisticsController
	taskId int
}

func (ts *unloadTaskStat) ApplyStat() {
	ts.sm.applyUnloadStat(ts.taskId)
}

///////////////////////////////////////////////////////////////////////////////

type collectTaskStat struct {
	sm           *StatisticsController
	taskId       int
	metricsCount int
	success      bool
	startTime    time.Time
	processTime  time.Time
}

func (ts *collectTaskStat) ApplyStat() {
	ts.sm.applyCollectStat(ts.taskId, ts.metricsCount, ts.success, ts.startTime, ts.processTime)
}
