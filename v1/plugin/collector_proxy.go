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

package plugin

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/solarwinds/snap-plugin-lib/v1/plugin/rpc"
)

const maxCollectChunkSize = 100

type collectorProxy struct {
	pluginProxy

	plugin Collector
}

func (c *collectorProxy) CollectMetrics(ctx context.Context, arg *rpc.MetricsArg) (*rpc.MetricsReply, error) {
	var logF = log.WithFields(log.Fields{"function": "CollectMetrics"})

	requestedMts := convertProtoToMetrics(arg.Metrics)

	collectedMts, err := c.plugin.CollectMetrics(requestedMts)
	if err != nil {
		return nil, err
	}

	protoMts, err := convertMetricsToProto(collectedMts)
	if err != nil {
		return nil, err
	}

	logF.WithFields(log.Fields{"length": len(arg.Metrics)}).Debug("Metrics will be sent to snap")
	return &rpc.MetricsReply{Metrics: protoMts}, nil
}

func (c *collectorProxy) CollectMetricsAsStream(arg *rpc.MetricsArg, stream rpc.Collector_CollectMetricsAsStreamServer) error {
	var logF = log.WithFields(log.Fields{"function": "CollectMetricsAsStream"})

	requestedMts := convertProtoToMetrics(arg.Metrics)
	collectedMts, err := c.plugin.CollectMetrics(requestedMts)
	if err != nil {
		return err
	}

	protoMts := make([]*rpc.Metric, 0, maxCollectChunkSize)
	for i, mt := range collectedMts {
		protoMt, err := toProtoMetric(mt)
		if err != nil {
			return err
		}
		protoMts = append(protoMts, protoMt)

		if len(protoMts) == maxCollectChunkSize || i == len(collectedMts)-1 {
			err := stream.Send(&rpc.MetricsReply{Metrics: protoMts})
			if err != nil {
				return err
			}
			logF.WithFields(log.Fields{"length": len(protoMts)}).Debug("Metrics chunk has been sent to snap")
			protoMts = make([]*rpc.Metric, 0, maxCollectChunkSize)
		}
	}

	return nil
}

func (c *collectorProxy) GetMetricTypes(ctx context.Context, arg *rpc.GetMetricTypesArg) (*rpc.MetricsReply, error) {
	cfg := fromProtoConfig(arg.Config)

	mts, err := c.plugin.GetMetricTypes(cfg)
	if err != nil {
		return nil, err
	}

	protoMts, err := convertMetricsToProto(mts)
	if err != nil {
		return nil, err
	}

	return &rpc.MetricsReply{Metrics: protoMts}, nil
}
