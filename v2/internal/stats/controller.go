package stats

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	statsChannelSize = 100
)

var log = logrus.WithFields(logrus.Fields{"layer": "lib", "module": "statistics"})

///////////////////////////////////////////////////////////////////////////////

type StatsController struct {
	startedSync     sync.Once
	incomingStatsCh chan StatCommand
	stats           Statistics
}

func NewStatsController(pluginName string, pluginVersion string) *StatsController {
	return &StatsController{
		startedSync: sync.Once{},

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

func (sc *StatsController) Run() {
	sc.startedSync.Do(func() {
		go func() {
			for {
				stat := <-sc.incomingStatsCh
				stat.ApplyStat()
			}
		}()
	})
}

func (sc *StatsController) UpdateLoadStat(taskId int, config string, filters []string) {
	sc.incomingStatsCh <- &loadTaskStat{
		sm:      sc,
		taskId:  taskId,
		config:  config,
		filters: filters,
	}
}

func (sc *StatsController) UpdateUnloadStat(taskId int) {
	sc.incomingStatsCh <- &unloadTaskStat{
		sm:     sc,
		taskId: taskId,
	}
}

func (sc *StatsController) UpdateCollectStat() {
	sc.incomingStatsCh <- &collectTaskStat{
		sm: sc,
	}
}

///////////////////////////////////////////////////////////////////////////////

func (sc *StatsController) applyLoadStat(taskId int, config string, filters []string) {
	// Update global stats
	sc.stats.tasks.CurrentlyActiveTasks += 1
	sc.stats.tasks.TotalActiveTasks += 1

	// Update task-specific stats
	sc.stats.tasksDetails[taskId] = taskDetailsFields{
		Configuration: config,
		Filters:       filters,
		LoadedTime:    time.Now(),
	}
}

func (sc *StatsController) applyUnloadStat(taskId int) {
	// Update global stats
	sc.stats.tasks.CurrentlyActiveTasks -= 1

	// Update task-specific stats
	delete(sc.stats.tasksDetails, taskId) // todo: safe?
}

func (sc *StatsController) applyCollectStat() {
	// Update global stats

	// Update task-specific state
}
