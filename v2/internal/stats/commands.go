package stats

///////////////////////////////////////////////////////////////////////////////

type loadTaskStat struct {
	sm *StatsManager
}

func (t *loadTaskStat) ApplyStat() {
	t.sm.applyLoadTaskStat()
}

///////////////////////////////////////////////////////////////////////////////

type unloadTaskStat struct {
	sm *StatsManager
}

func (t *unloadTaskStat) ApplyStat() {
	t.sm.applyUnloadTaskStat()
}

///////////////////////////////////////////////////////////////////////////////

type collectTaskStat struct {
	sm *StatsManager
}

func (t *collectTaskStat) ApplyStat() {
	t.sm.applyCollectStat()
}
