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

import (
	"time"
)

type MetricModifier interface {
	UpdateMetric(mt MetricSetter)
}

func MetricTag(key string, value string) MetricModifier {
	return &metricTags{
		tags: map[string]string{key: value},
	}
}

func MetricTags(tags map[string]string) MetricModifier {
	return &metricTags{
		tags: tags,
	}
}

func RemoveMetricTags(tags []string) MetricModifier {
	return &removeMetricTags{
		tagsToRemove: tags,
	}
}

func MetricTimestamp(timestamp time.Time) MetricModifier {
	return &metricTimestamp{
		timestamp: timestamp,
	}
}

func MetricDescription(description string) MetricModifier {
	return &metricDescription{
		description: description,
	}
}

func MetricUnit(unit string) MetricModifier {
	return &metricUnit{
		unit: unit,
	}
}

func MetricTypeGauge() MetricModifier {
	return &metricType{
		type_: GaugeType,
	}
}

func MetricTypeSum() MetricModifier {
	return &metricType{
		type_: SumType,
	}
}

func MetricTypeSummary() MetricModifier {
	return &metricType{
		type_: SummaryType,
	}
}

func MetricTypeHistogram() MetricModifier {
	return &metricType{
		type_: HistogramType,
	}
}

///////////////////////////////////////////////////////////////////////////////

type metricTags struct {
	tags map[string]string
}

func (m metricTags) UpdateMetric(mt MetricSetter) {
	mt.AddTags(m.tags)
}

type removeMetricTags struct {
	tagsToRemove []string
}

func (m removeMetricTags) UpdateMetric(mt MetricSetter) {
	mt.RemoveTags(m.tagsToRemove)
}

type metricTimestamp struct {
	timestamp time.Time
}

func (m metricTimestamp) UpdateMetric(mt MetricSetter) {
	mt.SetTimestamp(m.timestamp)
}

type metricDescription struct {
	description string
}

func (m metricDescription) UpdateMetric(mt MetricSetter) {
	mt.SetDescription(m.description)
}

type metricUnit struct {
	unit string
}

func (m metricUnit) UpdateMetric(mt MetricSetter) {
	mt.SetUnit(m.unit)
}

type metricType struct {
	type_ MetricType
}

func (m metricType) UpdateMetric(mt MetricSetter) {
	mt.SetType(m.type_)
}
