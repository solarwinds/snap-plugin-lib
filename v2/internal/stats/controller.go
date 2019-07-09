package stats

import (
	"sync"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
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
	closeCh         chan struct{}
	stats           Statistics
}

func NewStatsController(pluginName string, pluginVersion string) *StatsController {
	return &StatsController{
		startedSync:     sync.Once{},
		incomingStatsCh: make(chan StatCommand, statsChannelSize),
		closeCh:         make(chan struct{}),

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
				select {
				case stat := <-sc.incomingStatsCh:
					stat.ApplyStat()
				case <-sc.closeCh:
					return
				}
			}
		}()
	})
}

func (sc *StatsController) Close() {
	sc.closeCh <- struct{}{}
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

func (sc *StatsController) UpdateCollectStat(taskId int, mts []*types.Metric, success bool, startTime, endTime time.Time) {
	sc.incomingStatsCh <- &collectTaskStat{
		sm:          sc,
		taskId:      taskId,
		mts:         mts,
		success:     success,
		startTime:   startTime,
		processTime: endTime,
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

func (sc *StatsController) applyCollectStat(taskId int, mts []*types.Metric, success bool, startTime, completeTime time.Time) {
	duration := completeTime.Sub(startTime)

	// Update global stats
	sc.stats.tasks.TotalCollectsRequest += 1

	sc.stats.tasks.totalProcessingTime += duration
	sc.stats.tasks.AvgProcessingTime = 0 // todo
	sc.stats.tasks.MaxProcessingTime = 0 // todo

	// Update task-specific state
	taskStats := sc.stats.tasksDetails[taskId]

	taskStats.CollectRequest += 1
	taskStats.TotalMetrics += len(mts)
	taskStats.TotalProcessingTime += duration

	taskStats.AvgMetricsPerCollect = 0 // todo
	taskStats.AvgProcessingTime = 0    // todo
	taskStats.MaxProcessingTime = 0    // todo

	sc.stats.tasksDetails[taskId] = taskStats
}
