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

package service

import (
	"fmt"
	"time"

	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/types"
	"github.com/solarwinds/snap-plugin-lib/v2/pluginrpc"
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

	// when adding new type(s) apply changes also in PluginContext.isValidValueType() function
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
	case uint:
		grpcValue.DataVariant = &pluginrpc.MetricValue_VUint64{VUint64: uint64(t)}
	case []byte:
		grpcValue.DataVariant = &pluginrpc.MetricValue_VBytes{VBytes: t}
	case bool:
		grpcValue.DataVariant = &pluginrpc.MetricValue_VBool{VBool: t}
	case int16:
		grpcValue.DataVariant = &pluginrpc.MetricValue_VInt16{VInt16: int32(t)}
	case uint16:
		grpcValue.DataVariant = &pluginrpc.MetricValue_VUint16{VUint16: uint32(t)}
	case nil:
		grpcValue.DataVariant = nil
	default:
		return nil, fmt.Errorf("unsupported type: %T given in metric data", t)
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
	case *pluginrpc.MetricValue_VInt16:
		return int16(v.GetVInt16()), nil
	case *pluginrpc.MetricValue_VUint16:
		return uint16(v.GetVUint16()), nil
	}

	return nil, fmt.Errorf("unknown type of metric value: %T", v.DataVariant)
}

func toGRPCWarning(warning types.Warning) *pluginrpc.Warning {
	return &pluginrpc.Warning{
		Message:   warning.Message,
		Timestamp: toGRPCTime(warning.Timestamp),
	}
}
