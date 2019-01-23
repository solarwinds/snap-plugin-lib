package plugin

import (
	"github.com/librato/snap-plugin-lib-go/v1/plugin/rpc"
)

const DefaultMetricsChunkSize = 100

func convertMetricsToProto(mts []Metric) ([]*rpc.Metric, error) {
	var protoMts []*rpc.Metric

	for _, mt := range mts {
		protoMt, err := toProtoMetric(mt)
		if err != nil {
			return nil, err
		}
		protoMts = append(protoMts, protoMt)
	}

	return protoMts, nil
}

func convertProtoToMetrics(protoMts []*rpc.Metric) []Metric {
	var mts []Metric

	for _, protoMt := range protoMts {
		mt := fromProtoMetric(protoMt)
		mts = append(mts, mt)
	}

	return mts
}
