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

package service

import "github.com/solarwinds/snap-plugin-lib/v2/internal/util/types"

type CollectorProxy interface {
	RequestCollect(id string) <-chan types.CollectChunk
	LoadTask(id string, rawConfig []byte, mtsSelectors []string) error
	UnloadTask(id string) error
	CustomInfo(id string) ([]byte, error)
}
type PublisherProxy interface {
	RequestPublish(id string, mts []*types.Metric) types.ProcessingStatus
	LoadTask(id string, config []byte) error
	UnloadTask(id string) error
	CustomInfo(id string) ([]byte, error)
}
