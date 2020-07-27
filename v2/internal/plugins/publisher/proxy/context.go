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

package proxy

import (
	"errors"

	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

type PluginContext struct {
	*proxy.Context

	taskID     string
    sessionMts []*types.Metric
}

func NewPluginContext(ctxManager *ContextManager, taskID string, rawConfig []byte) (*PluginContext, error) {
	if ctxManager == nil {
		return nil, errors.New("can't create context without valid context manager")
	}

	baseContext, err := proxy.NewContext(rawConfig)
	if err != nil {
		return nil, err
	}

	return &PluginContext{
		Context: baseContext,
        taskID:  taskID,
	}, nil
}

func (pc *PluginContext) ListAllMetrics() []plugin.Metric {
	mts := make([]plugin.Metric, 0, len(pc.sessionMts))

	for _, mt := range pc.sessionMts {
		mts = append(mts, mt)
	}

	return mts
}

func (pc *PluginContext) Count() int {
	return len(pc.sessionMts)
}

func (pc *PluginContext) TaskID() string {
	return pc.taskID
}
