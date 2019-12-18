package pluginrpc

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/stats"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/pluginrpc"
)

// convert metric to GRPC structure
func toGRPCMetric(mt *types.Metric) (*pluginrpc.Metric, error) {
	value, err := toGRPCValue(mt.Value_)
	if err != nil {
		return nil, fmt.Errorf("can't convert metric to GRPC structure: %v", err)
	}

	protoMt := &pluginrpc.Metric{
		Namespace:   toGRPCNamespace(mt.Namespace_),
		Tags:        mt.Tags_,
		Value:       value,
		Unit:        mt.Unit_,
		Timestamp:   toGRPCTime(mt.Timestamp_),
		Description: mt.Description_,
	}

	return protoMt, nil
}

func fromGRPCMetric(mt *pluginrpc.Metric) (types.Metric, error) {
	data, err := fromGRPCValue(mt.Value)
	if err != nil {
		return types.Metric{}, fmt.Errorf("can't convert metric from GRPC structure: %v", err)
	}

	tags := map[string]string{}
	if mt.Tags != nil {
		tags = mt.Tags
	}

	retMt := types.Metric{
		Namespace_:   fromGRPCNamespace(mt.Namespace),
		Value_:       data,
		Tags_:        tags,
		Unit_:        mt.Unit,
		Timestamp_:   fromGRPCTime(mt.Timestamp),
		Description_: mt.Description,
	}

	return retMt, err
}

// convert namespace to GRPC structure
func toGRPCNamespace(ns []types.NamespaceElement) []*pluginrpc.Namespace {
	grpcNs := make([]*pluginrpc.Namespace, 0, len(ns))

	for _, nsElem := range ns {
		grpcNs = append(grpcNs, &pluginrpc.Namespace{
			Name:        nsElem.Name_,
			Value:       nsElem.Value_,
			Description: nsElem.Description_,
		})
	}

	return grpcNs
}

func fromGRPCNamespace(ns []*pluginrpc.Namespace) []types.NamespaceElement {
	retNsElem := make([]types.NamespaceElement, 0, len(ns))

	for _, nsElem := range ns {
		retNsElem = append(retNsElem, types.NamespaceElement{
			Name_:        nsElem.Name,
			Value_:       nsElem.Value,
			Description_: nsElem.Description,
		})
	}

	return retNsElem
}

func toGRPCTime(t time.Time) *pluginrpc.Time {
	return &pluginrpc.Time{
		Sec:  t.Unix(),
		Nsec: int64(t.Nanosecond()),
	}
}

func fromGRPCTime(t *pluginrpc.Time) time.Time {
	return time.Unix(t.Sec, t.Nsec)
}

// convert metric value to GRPC structure
func toGRPCValue(v interface{}) (*pluginrpc.MetricValue, error) {
	grpcValue := &pluginrpc.MetricValue{}

	switch t := v.(type) {
	case string:
		grpcValue.DataVariant = &pluginrpc.MetricValue_VString{VString: t}
	case float64:
		grpcValue.DataVariant = &pluginrpc.MetricValue_VDouble{VDouble: t}
	case float32:
		grpcValue.DataVariant = &pluginrpc.MetricValue_VFloat{VFloat: t}
	case int32:
		grpcValue.DataVariant = &pluginrpc.MetricValue_VInt32{VInt32: t}
	case int:
		grpcValue.DataVariant = &pluginrpc.MetricValue_VInt64{VInt64: int64(t)}
	case int64:
		grpcValue.DataVariant = &pluginrpc.MetricValue_VInt64{VInt64: t}
	case uint32:
		grpcValue.DataVariant = &pluginrpc.MetricValue_VUint32{VUint32: t}
	case uint64:
		grpcValue.DataVariant = &pluginrpc.MetricValue_VUint64{VUint64: t}
	case []byte:
		grpcValue.DataVariant = &pluginrpc.MetricValue_VBytes{VBytes: t}
	case bool:
		grpcValue.DataVariant = &pluginrpc.MetricValue_VBool{VBool: t}
	case nil:
		grpcValue.DataVariant = nil
	default:
		return nil, fmt.Errorf("unsupported type: %v given in metric data", t)
	}

	return grpcValue, nil
}

func fromGRPCValue(v *pluginrpc.MetricValue) (interface{}, error) {
	switch v.DataVariant.(type) {
	case *pluginrpc.MetricValue_VString:
		return v.GetVString(), nil
	case *pluginrpc.MetricValue_VDouble:
		return v.GetVDouble(), nil
	case *pluginrpc.MetricValue_VFloat:
		return v.GetVFloat(), nil
	case *pluginrpc.MetricValue_VInt32:
		return v.GetVInt32(), nil
	case *pluginrpc.MetricValue_VInt64:
		return v.GetVInt64(), nil
	case *pluginrpc.MetricValue_VUint32:
		return v.GetVInt32(), nil
	case *pluginrpc.MetricValue_VUint64:
		return v.GetVUint64(), nil
	case *pluginrpc.MetricValue_VBytes:
		return v.GetVBytes(), nil
	case *pluginrpc.MetricValue_VBool:
		return v.GetVBool(), nil
	}

	return nil, fmt.Errorf("unknown type of metric value: %T", v.DataVariant)
}

func toGRPCInfo(statistics *stats.Statistics, pprofLocation string) (*pluginrpc.Info, error) {
	pi := &statistics.PluginInfo
	ts := &statistics.TasksSummary

	info := &pluginrpc.Info{
		PluginInfo: &pluginrpc.PluginInfo{
			Name:           pi.Name,
			Version:        pi.Version,
			CmdLineOptions: pi.CmdLineOptions,
			Started:        toGRPCTime(pi.Started.Time),
		},
		TaskSummary: &pluginrpc.TaskSummary{
			Counters: &pluginrpc.TaskSummaryCounters{
				CurrentlyActiveTasks:   uint64(ts.Counters.CurrentlyActiveTasks),
				TotalActiveTasks:       uint64(ts.Counters.TotalActiveTasks),
				TotalExecutionRequests: uint64(ts.Counters.TotalExecutionRequests),
			},
			ProcessingTimes: &pluginrpc.ProcessingTimes{
				Total:   int64(ts.ProcessingTimes.Total),
				Average: int64(ts.ProcessingTimes.Average),
				Maximum: int64(ts.ProcessingTimes.Maximum),
			},
		},
		TaskDetails: map[string]*pluginrpc.TaskDetails{},
	}

	// Handle RawMessage - marshal to Json and unmarshal to typed struct
	b, err := pi.Options.MarshalJSON()
	if err != nil {
		return info, fmt.Errorf("could't marshal options field: %v", err)
	}

	options := &plugin.Options{}
	err = json.Unmarshal(b, options)
	if err != nil {
		return info, fmt.Errorf("could't unmarshal options field: %v", err)
	}

	info.PluginInfo.Options = &pluginrpc.Options{
		PluginIP:          options.PluginIP,
		GrpcPort:          uint32(options.GRPCPort),
		StatsPort:         uint32(options.StatsPort),
		GrpcPingTimeout:   int64(options.GRPCPingTimeout),
		GrpcPingMaxMissed: uint64(options.GRPCPingMaxMissed),
		LogLevel:          uint32(options.LogLevel),
		EnableProfiling:   options.EnableProfiling,
		ProfilingLocation: "",
		EnableStats:       options.EnableStats,
		EnableStatsServer: options.EnableStatsServer,
	}

	if options.EnableProfiling {
		info.PluginInfo.Options.ProfilingLocation = pprofLocation
	}

	for id, taskDetails := range statistics.TasksDetails {
		c := &taskDetails.Counters
		pt := &taskDetails.ProcessingTimes
		lm := &taskDetails.LastMeasurement

		info.TaskDetails[id] = &pluginrpc.TaskDetails{
			Configuration: fmt.Sprintf("%s", taskDetails.Configuration),
			Filters:       taskDetails.Filters,
			Counters: &pluginrpc.TaskDetailCounters{
				CollectRequests:            uint64(c.CollectRequests),
				TotalMetrics:               uint64(c.TotalMetrics),
				AverageMetricsPerExecution: uint64(c.AvgMetricsPerExecution),
			},
			Loaded: toGRPCTime(taskDetails.Loaded.Time),
			ProcessingTimes: &pluginrpc.ProcessingTimes{
				Total:   int64(pt.Total),
				Average: int64(pt.Average),
				Maximum: int64(pt.Maximum),
			},
			LastMeasurement: &pluginrpc.LastMeasurement{
				Occurred:         toGRPCTime(lm.Occurred.Time),
				ProcessedMetrics: uint64(lm.ProcessedMetrics),
				Duration:         int64(lm.Duration),
			},
		}
	}

	return info, nil
}
