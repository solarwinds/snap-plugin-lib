// +build medium

package plugin

import (
	"context"
	"errors"
)

type mockStreamer struct {
	mockPlugin
	err       error
	inMetric  chan []Metric
	outMetric chan []Metric
}

func newMockStreamer() *mockStreamer {
	return &mockStreamer{}
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
