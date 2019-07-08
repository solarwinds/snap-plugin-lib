package stats

type StatCommand interface {
	ApplyStat()
}

///////////////////////////////////////////////////////////////////////////////

type StatsManager struct {
	statsChan chan StatCommand
}

func NewStatsManager() *StatsManager {
	return &StatsManager{statsChan: make(chan StatCommand, 100)}
}

func (sm *StatsManager) Run() {
	for {
		stat := <-sm.statsChan
		stat.ApplyStat()
	}
}

func (sm *StatsManager) UpdateLoadTaskStat() {
	sm.statsChan <- &loadTaskStat{sm: sm}
}

func (sm *StatsManager) UpdateUnloadTaskStat() {
	sm.statsChan <- &unloadTaskStat{sm: sm}
}

func (sm *StatsManager) UpdateCollectStat() {
	sm.statsChan <- &collectTaskStat{sm: sm}
}

///////////////////////////////////////////////////////////////////////////////

func (sm *StatsManager) applyLoadTaskStat() {
	// todo: implement
}

func (sm *StatsManager) applyUnloadTaskStat() {
	// todo: implement
}

func (sm *StatsManager) applyCollectStat() {
	// todo: implement
}
