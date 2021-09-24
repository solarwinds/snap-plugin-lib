//go:build small
// +build small

/*
 Copyright 2016 Intel Corporation

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

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
	"context"
	"errors"
	"time"
)

func newMockPlugin() *mockPlugin {
	return &mockPlugin{}
}

func newMockErrPlugin() *mockPlugin {
	return &mockPlugin{err: errors.New("error")}
}

func newMockErrPublisher() *mockPublisher {
	return &mockPublisher{err: errors.New("error")}
}

func newMockErrCollector() *mockCollector {
	return &mockCollector{err: errors.New("empty")}
}

func newMockErrProcessor() *mockProcessor {
	return &mockProcessor{err: errors.New("error")}
}

type mockStreamer struct {
	mockPlugin
	err       error
	inMetric  chan []Metric
	outMetric chan []Metric
	action    func(chan []Metric, time.Duration, []Metric)
}

func newMockStreamer() *mockStreamer {
	return &mockStreamer{}
}

func newMockErrStreamer() *mockStreamer {
	return &mockStreamer{err: errors.New("empty")}
}

func newMockStreamerStream(action func(chan []Metric, time.Duration, []Metric)) *mockStreamer {
	return &mockStreamer{action: action}
}

func (mc *mockStreamer) doAction(t time.Duration, mts []Metric) {
	go func() {
		mc.action(mc.outMetric, t, mts)
	}()
}
func (mc *mockStreamer) GetMetricTypes(cfg Config) ([]Metric, error) {
	if mc.err != nil {
		return nil, errors.New("error")
	}

	mts := []Metric{}
	for _, v := range getMockMetricDataMap() {
		mts = append(mts, v)
	}
	return mts, nil
}

func (mc *mockStreamer) StreamMetrics(ctx context.Context, i chan []Metric, o chan []Metric, _ chan string) error {

	if mc.err != nil {
		return errors.New("error")
	}
	mc.inMetric = i
	mc.outMetric = o
	return nil
}
