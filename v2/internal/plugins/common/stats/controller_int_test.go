// +build medium

/*
 Copyright (c) 2020 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package stats

import (
	"testing"
	"time"

	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/types"
	"github.com/solarwinds/snap-plugin-lib/v2/plugin"

	stdCtx "context"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	pluginName    = "example"
	pluginVersion = "1.2.3"

	waitForCalculation = 200 * time.Millisecond
)

func TestStatistics(t *testing.T) {
	Convey("Validate that calculating statistics calculation is correct", t, func() {
		startTime := time.Unix(100000, 0)

		sci, _ := NewStatsController(stdCtx.Background(), pluginName, pluginVersion, types.PluginTypeCollector, &plugin.Options{})
		sc := sci.(*StatisticsController)

		// Load task1 and perform some collections
		{
			// Act
			sc.UpdateLoadStat("task-1", "cfg_1", []string{"filt_1_1", "filt_1_2", "filt_1_3"})
			sc.UpdateExecutionStat("task-1", 4, true, startTime.Add(1*time.Second), startTime.Add(3*time.Second))
			sc.UpdateExecutionStat("task-1", 6, true, startTime.Add(4*time.Second), startTime.Add(7*time.Second))
			sc.UpdateExecutionStat("task-1", 11, true, startTime.Add(8*time.Second), startTime.Add(12*time.Second))

			// Assert
			time.Sleep(waitForCalculation)

			pi := sc.stats.PluginInfo
			So(pi.Name, ShouldEqual, pluginName)
			So(pi.Version, ShouldEqual, pluginVersion)
			So(pi.Started.Time, ShouldNotBeNil)

			ts := sc.stats.TasksSummary
			So(ts.Counters.CurrentlyActiveTasks, ShouldEqual, 1)
			So(ts.Counters.TotalActiveTasks, ShouldEqual, 1)
			So(ts.Counters.TotalExecutionRequests, ShouldEqual, 3)

			td := sc.stats.TasksDetails
			So(td, ShouldContainKey, "task-1")
			So(td["task-1"].Counters.CollectRequests, ShouldEqual, 3)
			So(td["task-1"].Counters.TotalMetrics, ShouldEqual, 21)
			So(td["task-1"].LastMeasurement.ProcessedMetrics, ShouldEqual, 11)
			So(td["task-1"].ProcessingTimes.Total, ShouldEqual, 9*time.Second)
			So(td["task-1"].ProcessingTimes.Average, ShouldEqual, 3*time.Second)
		}

		// Load task2 and perform some collections
		{
			// Act
			sc.UpdateLoadStat("task-2", "cfg_1", []string{"filt_1_1", "filt_1_2", "filt_1_3"})
			sc.UpdateExecutionStat("task-2", 5, true, startTime.Add(20*time.Second), startTime.Add(21*time.Second))
			sc.UpdateExecutionStat("task-2", 15, true, startTime.Add(25*time.Second), startTime.Add(26*time.Second))
			sc.UpdateExecutionStat("task-2", 10, true, startTime.Add(30*time.Second), startTime.Add(34*time.Second))

			// Assert
			time.Sleep(waitForCalculation)

			ts := sc.stats.TasksSummary
			So(ts.Counters.CurrentlyActiveTasks, ShouldEqual, 2)
			So(ts.Counters.TotalActiveTasks, ShouldEqual, 2)
			So(ts.Counters.TotalExecutionRequests, ShouldEqual, 6)

			td := sc.stats.TasksDetails
			So(td, ShouldContainKey, "task-1")
			So(td, ShouldContainKey, "task-2")
			So(td["task-2"].Counters.CollectRequests, ShouldEqual, 3)
			So(td["task-2"].Counters.TotalMetrics, ShouldEqual, 30)
			So(td["task-2"].LastMeasurement.ProcessedMetrics, ShouldEqual, 10)
			So(td["task-2"].ProcessingTimes.Total, ShouldEqual, 6*time.Second)
			So(td["task-2"].ProcessingTimes.Average, ShouldEqual, 2*time.Second)
		}

		// Unload task1
		{
			// Act
			sc.UpdateUnloadStat("task-1")

			// Assert
			time.Sleep(waitForCalculation)

			ts := sc.stats.TasksSummary
			So(ts.Counters.CurrentlyActiveTasks, ShouldEqual, 1)
			So(ts.Counters.TotalActiveTasks, ShouldEqual, 2)
			So(ts.Counters.TotalExecutionRequests, ShouldEqual, 6)

			td := sc.stats.TasksDetails
			So(td, ShouldNotContainKey, "task-1")
			So(td, ShouldContainKey, "task-2")
		}

		// Load task3 and perform some operations
		{
			// Act
			sc.UpdateLoadStat("task-3", "cfg_1", []string{"filt_1_1", "filt_1_2", "filt_1_3"})

			sc.UpdateExecutionStat("task-3", 1, true, startTime.Add(40*time.Second), startTime.Add(41*time.Second))
			sc.UpdateExecutionStat("task-3", 0, true, startTime.Add(45*time.Second), startTime.Add(46*time.Second))

			sc.UpdateExecutionStat("task-2", 3, true, startTime.Add(50*time.Second), startTime.Add(51*time.Second))

			// Assert
			time.Sleep(waitForCalculation)

			ts := sc.stats.TasksSummary
			So(ts.Counters.CurrentlyActiveTasks, ShouldEqual, 2)
			So(ts.Counters.TotalActiveTasks, ShouldEqual, 3)
			So(ts.Counters.TotalExecutionRequests, ShouldEqual, 9)

			td := sc.stats.TasksDetails
			So(td, ShouldContainKey, "task-2")
			So(td, ShouldContainKey, "task-3")

			So(td["task-2"].Counters.CollectRequests, ShouldEqual, 4)
			So(td["task-2"].Counters.TotalMetrics, ShouldEqual, 33)
			So(td["task-2"].LastMeasurement.ProcessedMetrics, ShouldEqual, 3)

			So(td["task-3"].Counters.CollectRequests, ShouldEqual, 2)
			So(td["task-3"].Counters.TotalMetrics, ShouldEqual, 1)
			So(td["task-3"].LastMeasurement.ProcessedMetrics, ShouldEqual, 0)
		}

		// Unload task2 and task3
		{
			// Act
			sc.UpdateUnloadStat("task-2")
			sc.UpdateUnloadStat("task-3")

			// Assert
			time.Sleep(waitForCalculation)

			ts := sc.stats.TasksSummary
			So(ts.Counters.CurrentlyActiveTasks, ShouldEqual, 0)
			So(ts.Counters.TotalActiveTasks, ShouldEqual, 3)
			So(ts.Counters.TotalExecutionRequests, ShouldEqual, 9)

			td := sc.stats.TasksDetails
			So(td, ShouldNotContainKey, "task-1")
			So(td, ShouldNotContainKey, "task-2")
			So(td, ShouldNotContainKey, "task-3")
		}

		// Finalize
		sc.Close()
	})
}
