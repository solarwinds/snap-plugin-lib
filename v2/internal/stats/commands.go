package stats

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
	sm *StatsController
}

func (ts *collectTaskStat) ApplyStat() {
	ts.sm.applyCollectStat()
}
