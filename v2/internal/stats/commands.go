package stats

import (
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
)

///////////////////////////////////////////////////////////////////////////////

type StatCommand interface {
	ApplyStat()
}

///////////////////////////////////////////////////////////////////////////////

type loadTaskStat struct {
	sm      *StatsController
	taskId  int
	config  string
	filters []string
}

func (ts *loadTaskStat) ApplyStat() {
	ts.sm.applyLoadStat(ts.taskId, ts.config, ts.filters)
}

///////////////////////////////////////////////////////////////////////////////

type unloadTaskStat struct {
	sm     *StatsController
	taskId int
}

func (ts *unloadTaskStat) ApplyStat() {
	ts.sm.applyUnloadStat(ts.taskId)
}

///////////////////////////////////////////////////////////////////////////////

type collectTaskStat struct {
	sm          *StatsController
	taskId      int
	mts         []*types.Metric
	success     bool
	startTime   time.Time
	processTime time.Time
}

func (ts *collectTaskStat) ApplyStat() {
	ts.sm.applyCollectStat(ts.taskId, ts.mts, ts.success, ts.startTime, ts.processTime)
}
