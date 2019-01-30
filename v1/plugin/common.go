package plugin

import (
	"github.com/librato/snap-plugin-lib-go/v1/plugin/rpc"
)

func convertMetricsToProto(mts []Metric) ([]*rpc.Metric, error) {
	protoMts := make([]*rpc.Metric, 0, len(mts))

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
	mts := make([]Metric, 0, len(protoMts))

	for _, protoMt := range protoMts {
		mt := fromProtoMetric(protoMt)
		mts = append(mts, mt)
	}

	return mts
}
