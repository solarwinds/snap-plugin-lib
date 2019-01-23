/*
http://www.apache.org/licenses/LICENSE-2.0.txt


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
*/

package plugin

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/librato/snap-plugin-lib-go/v1/plugin/rpc"
)

// TODO(danielscottt): plugin panics

type collectorProxy struct {
	pluginProxy

	plugin Collector
}

func (c *collectorProxy) CollectMetrics(ctx context.Context, arg *rpc.MetricsArg) (*rpc.MetricsReply, error) {
	requestedMts := convertProtoToMetrics(arg.Metrics)

	collectedMts, err := c.plugin.CollectMetrics(requestedMts)
	if err != nil {
		return nil, err
	}

	protoMts, err := convertMetricsToProto(collectedMts)
	if err != nil {
		return nil, err
	}

	return &rpc.MetricsReply{Metrics: protoMts}, nil
}

func (c *collectorProxy) CollectMetricsAsStream(arg *rpc.MetricsArg, stream rpc.Collector_CollectMetricsAsStreamServer) error {
	logF := logrus.WithFields(logrus.Fields{"test": "adamiklib", "block": "CollectMetricsAsStream"})

	requestedMts := convertProtoToMetrics(arg.Metrics)

	collectedMts, err := c.plugin.CollectMetrics(requestedMts)
	if err != nil {
		return err
	}

	splitMts := ChunkMetrics(collectedMts, DefaultMetricsChunkSize)

	logF.Debugf("There are %d metrics to send. They might be send in chunks", len(collectedMts))
	for _, chunkMts := range splitMts {
		protoMts, err := convertMetricsToProto(chunkMts)
		if err != nil {
			return err
		}

		stream.Send(&rpc.MetricsReply{Metrics: protoMts})
		logF.Debugf("Chunk sent to snap (len=%d)", len(protoMts))
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
