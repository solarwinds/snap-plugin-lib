package pluginrpc

import (
	"fmt"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/stats"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
)

// convert metric to GRPC structure
func toGRPCMetric(mt *types.Metric) (*Metric, error) {
	value, err := toGRPCValue(mt.Value)
	if err != nil {
		return nil, fmt.Errorf("can't convert metric to GRPC structure: %v", err)
	}

	protoMt := &Metric{
		Namespace: toGRPCNamespace(mt.Namespace),
		Tags:      mt.Tags,
		Value:     value,
		Unit:      mt.Unit,
		Timestamp: toGRPCTime(mt.Timestamp),
	}

	return protoMt, nil
}

// convert namespace to GRPC structure
func toGRPCNamespace(ns []types.NamespaceElement) []*Namespace {
	grpcNs := make([]*Namespace, 0, len(ns))

	for _, nsElem := range ns {
		grpcNs = append(grpcNs, &Namespace{
			Name:        nsElem.Name,
			Value:       nsElem.Value,
			Description: nsElem.Description,
		})
	}

	return grpcNs
}

func toGRPCTime(t time.Time) *Time {
	return &Time{
		Sec:  t.Unix(),
		Nsec: int64(t.Nanosecond()),
	}
}

// convert metric value to GRPC structure
func toGRPCValue(v interface{}) (*MetricValue, error) {
	grpcValue := &MetricValue{}

	switch t := v.(type) {
	case string:
		grpcValue.DataVariant = &MetricValue_VString{VString: t}
	case float64:
		grpcValue.DataVariant = &MetricValue_VDouble{VDouble: t}
	case float32:
		grpcValue.DataVariant = &MetricValue_VFloat{VFloat: t}
	case int32:
		grpcValue.DataVariant = &MetricValue_VInt32{VInt32: t}
	case int:
		grpcValue.DataVariant = &MetricValue_VInt64{VInt64: int64(t)}
	case int64:
		grpcValue.DataVariant = &MetricValue_VInt64{VInt64: t}
	case uint32:
		grpcValue.DataVariant = &MetricValue_VUint32{VUint32: t}
	case uint64:
		grpcValue.DataVariant = &MetricValue_VUint64{VUint64: t}
	case []byte:
		grpcValue.DataVariant = &MetricValue_VBytes{VBytes: t}
	case bool:
		grpcValue.DataVariant = &MetricValue_VBool{VBool: t}
	case nil:
		grpcValue.DataVariant = nil
	default:
		return nil, fmt.Errorf("unsupported type: %v given in metric data", t)
	}

	return grpcValue, nil
}

func toGRPCInfo(statistics stats.Statistics) *Info {
	pi := &statistics.PluginInfo
	ts := &statistics.TasksSummary

	info := &Info{
		PluginInfo: &PluginInfo{
			Name:           pi.Name,
			Version:        pi.Version,
			CmdLineOptions: pi.CmdLineOptions,
			Options:        nil, // todo: complete
			Started:        toGRPCTime(pi.Started.Time),
		},
		TaskSummary: &TaskSummary{
			Counters: &TaskSummaryCounters{
				CurrentlyActiveTasks: uint64(ts.Counters.CurrentlyActiveTasks), // todo: uint64 types in original code
				TotalActiveTasks:     uint64(ts.Counters.TotalActiveTasks),
				TotalCollectRequests: uint64(ts.Counters.TotalCollectRequests),
			},
			ProcessingTimes: &ProcessingTimes{
				Total:   int64(ts.ProcessingTimes.Total),
				Average: int64(ts.ProcessingTimes.Average),
				Maximum: int64(ts.ProcessingTimes.Maximum),
			},
		},
		TaskDetails: map[uint64]*TaskDetails{},
	}

	for id, taskDetails := range statistics.TasksDetails {
		c := &taskDetails.Counters
		pt := &taskDetails.ProcessingTimes
		lm := &taskDetails.LastMeasurement

		info.TaskDetails[uint64(id)] = &TaskDetails{ // todo: conversion
			Configuration: "", // todo
			Filters:       taskDetails.Filters,
			Counters: &TaskDetailCounters{
				CollectRequests:          uint64(c.CollectRequests),
				TotalMetrics:             uint64(c.TotalMetrics),
				AverageMetricsPerCollect: uint64(c.AvgMetricsPerCollect),
			},
			Loaded: toGRPCTime(taskDetails.Loaded.Time),
			ProcessingTimes: &ProcessingTimes{
				Total:   int64(pt.Total),
				Average: int64(pt.Average),
				Maximum: int64(pt.Maximum),
			},
			LastMeasurement: &LastMeasurement{
				Occurred:         toGRPCTime(lm.Occurred.Time),
				CollectedMetrics: uint64(lm.CollectedMetrics),
			},
		}
	}

	return info
}
