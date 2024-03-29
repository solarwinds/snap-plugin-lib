/*
 Copyright (c) 2021 SolarWinds Worldwide, LLC

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

package plugin

type MetricType int

const (
	UnknownType MetricType = iota
	GaugeType
	SumType
	SummaryType
	HistogramType
)

type Histogram struct {
	DataPoints map[float64]float64
	Count      int
	Sum        float64
}

// TODO: https://swicloud.atlassian.net/browse/AO-20547: Add support for Summary Quantiles
type Summary struct {
	Count int
	Sum   float64
}
