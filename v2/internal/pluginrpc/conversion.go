package pluginrpc

import (
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

func toGRPCMetric(mt *plugin.Metric) (*Metric, error) {
	protoMt := &Metric{
		Namespace: mt.Namespace,
		Tags:      mt.Tags,
		Value:     nil, // todo: fix sending values
		Timestamp: &Time{
			Sec:  mt.Timestamp.Unix(),
			Nsec: int64(mt.Timestamp.Nanosecond()),
		},
	}

	return protoMt, nil
}
