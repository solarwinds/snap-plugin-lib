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

package mock

import "github.com/stretchr/testify/mock"

type Definition struct {
	mock.Mock
}

func (m *Definition) DefineTasksPerInstanceLimit(limit int) error {
	args := m.Called(limit)
	return args.Error(0)
}

func (m *Definition) DefineInstancesLimit(limit int) error {
	args := m.Called(limit)
	return args.Error(0)
}

type CollectorDefinition struct {
	Definition
}

func (m *CollectorDefinition) DefineMetric(namespace string, unit string, isDefault bool, description string) {
	m.Called(namespace, unit, isDefault, description)
}

func (m *CollectorDefinition) DefineGroup(name string, description string) {
	m.Called(name, description)
}

func (m *CollectorDefinition) SetAllowDynamicLastElement() {
	m.Called()
}

func (m *CollectorDefinition) SetAllowAddingUndefinedMetrics() {
	m.Called()
}

func (m *CollectorDefinition) SetAllowValuesAtAnyNamespaceLevel() {
	m.Called()
}

func (m *CollectorDefinition) DefineExampleConfig(cfg string) error {
	args := m.Called(cfg)
	return args.Error(0)
}
