package stats

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"

	"github.com/sirupsen/logrus"
)

const (
	statsChannelSize = 100
	reqChannelSize   = 10
)

var log = logrus.WithFields(logrus.Fields{"layer": "lib", "module": "statistics"})

///////////////////////////////////////////////////////////////////////////////

type Controller interface {
	Close()
	RequestStat() chan *Statistics
	UpdateLoadStat(taskID string, config string, filters []string)
	UpdateUnloadStat(taskID string)
	UpdateExecutionStat(taskID string, metricsCount int, success bool, startTime, endTime time.Time)
	UpdateStreamingStat(taskID string, metricsCount int, startTime, lastUpdate time.Time)
}

///////////////////////////////////////////////////////////////////////////////

func NewController(pluginName string, pluginVersion string, pluginType types.PluginType, opt *plugin.Options) (Controller, error) {
	if opt.EnableStats {
		return NewStatsController(pluginName, pluginVersion, pluginType, opt)
	}

	return NewEmptyController()
}

///////////////////////////////////////////////////////////////////////////////

type StatisticsController struct {
	pluginType types.PluginType

	startedSync       sync.Once
	incomingStatsCh   chan StatCommand
	incomingRequestCh chan chan *Statistics
	closeCh           chan struct{}
	stats             *Statistics
}

func NewStatsController(pluginName string, pluginVersion string, pluginType types.PluginType, opt *plugin.Options) (Controller, error) {
	optJson, err := json.Marshal(opt)
	if err != nil {
		return nil, err
	}

	sc := &StatisticsController{
		pluginType: pluginType,

		startedSync:       sync.Once{},
		incomingStatsCh:   make(chan StatCommand, statsChannelSize),
		incomingRequestCh: make(chan chan *Statistics, reqChannelSize),
		closeCh:           make(chan struct{}),

		stats: &Statistics{
			PluginInfo: pluginInfo{
				Name:    pluginName,
				Version: pluginVersion,
				Type:    string(pluginType),
				Started: eventTimes{
					Time: time.Now(),
				},

				CmdLineOptions: strings.Join(os.Args[1:], " "),
				Options:        optJson,
			},
			TasksSummary: tasksSummary{},
			TasksDetails: map[string]taskDetails{},
		},
	}

	sc.run()

	return sc, nil
}

func (sc *StatisticsController) run() {
	sc.startedSync.Do(func() {
		go func() {
			for {
				select {
				case stat := <-sc.incomingStatsCh:
					stat.ApplyStat()
				case respCh := <-sc.incomingRequestCh:
					respCh <- sc.stats
				case <-sc.closeCh:
					return
				}
			}
		}()
	})
}

func (sc *StatisticsController) Close() {
	sc.closeCh <- struct{}{}
}

func (sc *StatisticsController) RequestStat() chan *Statistics {
	respCh := make(chan *Statistics)

	sc.incomingRequestCh <- respCh

	return respCh
}

func (sc *StatisticsController) UpdateLoadStat(taskID string, config string, filters []string) {
	sc.incomingStatsCh <- &loadTaskStat{
		sm:      sc,
		taskID:  taskID,
		config:  config,
		filters: filters,
	}
}

func (sc *StatisticsController) UpdateUnloadStat(taskID string) {
	sc.incomingStatsCh <- &unloadTaskStat{
		sm:     sc,
		taskID: taskID,
	}
}

func (sc *StatisticsController) UpdateExecutionStat(taskID string, metricsCount int, success bool, startTime, endTime time.Time) {
	sc.incomingStatsCh <- &collectTaskStat{
		sm:           sc,
		taskID:       taskID,
		metricsCount: metricsCount,
		success:      success,
		startTime:    startTime,
		processTime:  endTime,
	}
}

func (sc *StatisticsController) UpdateStreamingStat(taskID string, metricsCount int, startTime, lastUpdate time.Time) {
	sc.incomingStatsCh <- &streamTaskStat{
		sm:           sc,
		taskID:       taskID,
		metricsCount: metricsCount,
		startTime:    startTime,
		lastUpdate:   lastUpdate,
	}
}

///////////////////////////////////////////////////////////////////////////////

func (sc *StatisticsController) applyLoadStat(taskID string, config string, filters []string) {
	log.WithFields(logrus.Fields{
		"task-id":        taskID,
		"statistic-type": "Load",
	}).Trace("Applying statistic")

	// Update global stats
	sc.stats.TasksSummary.Counters.CurrentlyActiveTasks += 1
	sc.stats.TasksSummary.Counters.TotalActiveTasks += 1

	if filters == nil { // generate [] instead of null when marshaling
		filters = []string{}
	}

	// Update task-specific stats
	sc.stats.TasksDetails[taskID] = taskDetails{
		Configuration: json.RawMessage(config),
		Filters:       filters,
		Loaded: eventTimes{
			Time: time.Now(),
		},
	}
}

func (sc *StatisticsController) applyUnloadStat(taskID string) {
	log.WithFields(logrus.Fields{
		"task-id":        taskID,
		"statistic-type": "Unload",
	}).Trace("Applying statistic")

	// Update global stats
	sc.stats.TasksSummary.Counters.CurrentlyActiveTasks -= 1

	// Update task-specific stats
	delete(sc.stats.TasksDetails, taskID)
}

func (sc *StatisticsController) applyCollectStat(taskID string, metricsCount int, _ bool, startTime, completeTime time.Time) {
	log.WithFields(logrus.Fields{
		"task-id":        taskID,
		"statistic-type": "Collect",
	}).Trace("Applying statistic")
	processingTime := completeTime.Sub(startTime)

	// Update global stats
	{
		ts := &sc.stats.TasksSummary

		ts.Counters.TotalExecutionRequests += 1
		ts.ProcessingTimes.Total += processingTime

		if ts.Counters.TotalExecutionRequests > 0 {
			ts.ProcessingTimes.Average = time.Duration(int(ts.ProcessingTimes.Total) / ts.Counters.TotalExecutionRequests)
		}

		if processingTime > ts.ProcessingTimes.Maximum {
			ts.ProcessingTimes.Maximum = processingTime
		}
	}

	// Update task-specific state
	{
		td := sc.stats.TasksDetails[taskID]

		td.Counters.CollectRequests += 1
		td.Counters.TotalMetrics += metricsCount
		td.ProcessingTimes.Total += processingTime

		if td.Counters.CollectRequests > 0 {
			td.ProcessingTimes.Average = time.Duration(int(td.ProcessingTimes.Total) / td.Counters.CollectRequests)
			td.Counters.AvgMetricsPerExecution = td.Counters.TotalMetrics / td.Counters.CollectRequests
		}

		if processingTime > td.ProcessingTimes.Maximum {
			td.ProcessingTimes.Maximum = processingTime
		}

		td.LastMeasurement = measurementInfo{
			Occurred: eventTimes{
				Time: completeTime,
			},
			Duration:         processingTime,
			ProcessedMetrics: metricsCount,
		}

		sc.stats.TasksDetails[taskID] = td
	}
}

func (sc *StatisticsController) applyStreamStat(taskID string, metricsCount int, startTime, lastUpdate time.Time) {
	log.WithFields(logrus.Fields{
		"task-id":        taskID,
		"statistic-type": "Streaming",
	}).Trace("Applying statistic")
	processingTime := lastUpdate.Sub(startTime)

	td := sc.stats.TasksDetails[taskID]
	td.ProcessingTimes.Total = processingTime
	td.Counters.CollectRequests = 1
	td.Counters.TotalMetrics += metricsCount

	sc.stats.TasksDetails[taskID] = td
}

///////////////////////////////////////////////////////////////////////////////

func NewEmptyController() (Controller, error) {
	return &EmptyController{}, nil
}

type EmptyController struct {
}

func (d *EmptyController) Close() {
}

func (d *EmptyController) RequestStat() chan *Statistics {
	statCh := make(chan *Statistics)

	go func() {
		statCh <- nil
	}()

	return statCh
}

func (d *EmptyController) UpdateLoadStat(taskID string, config string, filters []string) {
}

func (d *EmptyController) UpdateUnloadStat(taskID string) {
}

func (d *EmptyController) UpdateExecutionStat(taskID string, metricsCount int, success bool, startTime, endTime time.Time) {
}

func (d *EmptyController) UpdateStreamingStat(taskID string, metricsCount int, startTime, lastUpdate time.Time) {
}
