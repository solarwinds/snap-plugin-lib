package stats

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/sirupsen/logrus"
)

const (
	statsChannelSize = 100
	reqChannelSize   = 10
)

var log = logrus.WithFields(logrus.Fields{"layer": "lib", "module": "statistics"})

///////////////////////////////////////////////////////////////////////////////

type StatsController struct {
	startedSync       sync.Once
	incomingStatsCh   chan StatCommand
	incomingRequestCh chan chan Statistics
	closeCh           chan struct{}
	stats             Statistics
}

func NewStatsController(pluginName string, pluginVersion string) *StatsController {
	return &StatsController{
		startedSync:       sync.Once{},
		incomingStatsCh:   make(chan StatCommand, statsChannelSize),
		incomingRequestCh: make(chan chan Statistics, reqChannelSize),
		closeCh:           make(chan struct{}),

		stats: Statistics{
			PluginInfo: pluginInfoFields{
				Name:      pluginName,
				Version:   pluginVersion,
				StartTime: time.Now(),
			},
			Tasks:        tasksFields{},
			TasksDetails: map[int]taskDetailsFields{},
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
				case respCh := <-sc.incomingRequestCh:
					sc.refresh()
					respCh <- sc.stats
					close(respCh)
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

func (sc *StatsController) RequestStat() chan Statistics {
	respCh := make(chan Statistics)

	sc.incomingRequestCh <- respCh

	return respCh
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
	sc.stats.Tasks.CurrentlyActiveTasks += 1
	sc.stats.Tasks.TotalActiveTasks += 1

	// Update task-specific stats
	sc.stats.TasksDetails[taskId] = taskDetailsFields{
		Configuration: json.RawMessage(config),
		Filters:       filters,
		LoadedTime:    time.Now(),
	}
}

func (sc *StatsController) applyUnloadStat(taskId int) {
	// Update global stats
	sc.stats.Tasks.CurrentlyActiveTasks -= 1

	// Update task-specific stats
	delete(sc.stats.TasksDetails, taskId) // todo: safe?
}

func (sc *StatsController) applyCollectStat(taskId int, mts []*types.Metric, success bool, startTime, completeTime time.Time) {
	processingTime := completeTime.Sub(startTime)

	// Update global stats
	{
		sc.stats.Tasks.TotalCollectsRequest += 1

		sc.stats.Tasks.totalProcessingTime += processingTime
		sc.stats.Tasks.AvgProcessingTime = 0 // todo
		sc.stats.Tasks.MaxProcessingTime = 0 // todo
	}

	// Update task-specific state
	{
		taskStats := sc.stats.TasksDetails[taskId]

		taskStats.CollectRequest += 1
		taskStats.TotalMetrics += len(mts)
		taskStats.totalProcessingTime += processingTime

		if taskStats.CollectRequest > 0 {
			taskStats.avgProcessingTime = time.Duration(int(taskStats.totalProcessingTime) / taskStats.CollectRequest)
			taskStats.AvgMetricsPerCollect = taskStats.TotalMetrics / taskStats.CollectRequest
		}
		if processingTime > taskStats.maxProcessingTime {
			taskStats.maxProcessingTime = processingTime
		}

		sc.stats.TasksDetails[taskId] = taskStats
	}

}

func (sc *StatsController) refresh() {
	for id, taskStat := range sc.stats.TasksDetails {
		taskStat.TotalProcessingTime = taskStat.totalProcessingTime.String()
		taskStat.AvgProcessingTime = taskStat.avgProcessingTime.String()
		taskStat.MaxProcessingTime = taskStat.maxProcessingTime.String()

		sc.stats.TasksDetails[id] = taskStat
	}
}
