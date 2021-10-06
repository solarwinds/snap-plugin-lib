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

package types

import (
	"fmt"
	"time"

	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
)

const (
	metricSeparator = "/"
)

type Metric struct {
	Namespace_   []NamespaceElement
	Value_       interface{}
	Tags_        map[string]string
	Unit_        string
	Timestamp_   time.Time
	Description_ string
	Type_        plugin.MetricType
}

func (m Metric) Namespace() plugin.Namespace {
	ns := make(Namespace, 0, len(m.Namespace_))

	for i := range m.Namespace_ {
		ns = append(ns, m.Namespace_[i])
	}

	return ns
}

func (m Metric) Value() interface{} {
	return m.Value_
}

func (m Metric) Tags() map[string]string {
	return m.Tags_
}

func (m Metric) Unit() string {
	return m.Unit_
}

func (m Metric) Description() string {
	return m.Description_
}

func (m Metric) Type() plugin.MetricType {
	return m.Type_
}

func (m Metric) Timestamp() time.Time {
	return m.Timestamp_
}

func (m Metric) String() string {
	return fmt.Sprintf("%s %v {%v}", m.Namespace().String(), m.Value_, m.Tags_)
}

func (m *Metric) AddTags(tags map[string]string) {
	if m.Tags_ == nil { // lazy initialization
		m.Tags_ = map[string]string{}
	}

	for k, v := range tags {
		m.Tags_[k] = v
	}
}

func (m *Metric) RemoveTags(keys []string) {
	for _, k := range keys {
		delete(m.Tags_, k)
	}
}

func (m *Metric) SetDescription(description string) {
	m.Description_ = description
}

func (m *Metric) SetUnit(unit string) {
	m.Unit_ = unit
}

func (m *Metric) SetTimestamp(timestamp time.Time) {
	m.Timestamp_ = timestamp
}

func (m *Metric) SetType(type_ plugin.MetricType) {
	m.Type_ = type_
}
