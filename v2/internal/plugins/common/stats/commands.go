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
	"time"
)

///////////////////////////////////////////////////////////////////////////////

type StatCommand interface {
	ApplyStat()
}

///////////////////////////////////////////////////////////////////////////////

type loadTaskStat struct {
	sm      *StatisticsController
	taskID  string
	config  string
	filters []string
}

func (ts *loadTaskStat) ApplyStat() {
	ts.sm.applyLoadStat(ts.taskID, ts.config, ts.filters)
}

///////////////////////////////////////////////////////////////////////////////

type unloadTaskStat struct {
	sm     *StatisticsController
	taskID string
}

func (ts *unloadTaskStat) ApplyStat() {
	ts.sm.applyUnloadStat(ts.taskID)
}

///////////////////////////////////////////////////////////////////////////////

type collectTaskStat struct {
	sm           *StatisticsController
	taskID       string
	metricsCount int
	success      bool
	startTime    time.Time
	processTime  time.Time
}

func (ts *collectTaskStat) ApplyStat() {
	ts.sm.applyCollectStat(ts.taskID, ts.metricsCount, ts.success, ts.startTime, ts.processTime)
}

///////////////////////////////////////////////////////////////////////////////

type streamTaskStat struct {
	sm           *StatisticsController
	taskID       string
	metricsCount int
	startTime    time.Time
	lastUpdate   time.Time
}

func (ts *streamTaskStat) ApplyStat() {
	ts.sm.applyStreamStat(ts.taskID, ts.metricsCount, ts.startTime, ts.lastUpdate)
}
