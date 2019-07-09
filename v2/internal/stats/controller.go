package stats

import (
	"time"

	"github.com/sirupsen/logrus"
)

const (
	statsChannelSize = 100
)

var log = logrus.WithFields(logrus.Fields{"layer": "lib", "module": "statistics"})

///////////////////////////////////////////////////////////////////////////////

type StatsController struct {
	incomingStatsCh chan StatCommand
	stats           Statistics
}

func NewStatsController(pluginName string, pluginVersion string) *StatsController {
	return &StatsController{
		incomingStatsCh: make(chan StatCommand, statsChannelSize),
		stats: Statistics{
			pluginInfo: pluginInfoFields{
				Name:      pluginName,
				Version:   pluginVersion,
				StartTime: time.Now(),
			},
			tasks: tasksFields{},
		},
	}
}

func (sm *StatsController) Run() {
	for {
		stat := <-sm.incomingStatsCh
		stat.ApplyStat()
	}
}

func (sm *StatsController) UpdateLoadTaskStat() {
	sm.incomingStatsCh <- &loadTaskStat{sm: sm}
}

func (sm *StatsController) UpdateUnloadTaskStat() {
	sm.incomingStatsCh <- &unloadTaskStat{sm: sm}
}

func (sm *StatsController) UpdateCollectStat() {
	sm.incomingStatsCh <- &collectTaskStat{sm: sm}
}

///////////////////////////////////////////////////////////////////////////////

func (sm *StatsController) applyLoadTaskStat() {
	// todo: implement
}

func (sm *StatsController) applyUnloadTaskStat() {
	// todo: implement
}

func (sm *StatsController) applyCollectStat() {
	// todo: implement
}
