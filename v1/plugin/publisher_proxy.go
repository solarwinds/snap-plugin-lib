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
	"context"
	"fmt"
	"io"

	"github.com/librato/snap-plugin-lib-go/v1/plugin/rpc"
	log "github.com/sirupsen/logrus"
)

type publisherProxy struct {
	pluginProxy

	plugin Publisher
}

func (p *publisherProxy) Publish(ctx context.Context, arg *rpc.PubProcArg) (*rpc.ErrReply, error) {
	var logF = log.WithFields(log.Fields{"function": "Publish"})

	mts := convertProtoToMetrics(arg.Metrics)
	cfg := fromProtoConfig(arg.Config)

	logF.WithFields(log.Fields{"length": len(mts)}).Debug("Metrics will be sent to appoptics")
	err := p.plugin.Publish(mts, cfg)
	if err != nil {
		return &rpc.ErrReply{Error: err.Error()}, nil
	}
	return &rpc.ErrReply{}, nil
}

func (p *publisherProxy) PublishAsStream(stream rpc.Publisher_PublishAsStreamServer) error {
	var logF = log.WithFields(log.Fields{"function": "PublishAsStream"})

	var mts []Metric
	var cfg Config

	for {
		protoMts, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			return fmt.Errorf("failure when reading from stream: %s", err.Error())
		}

		logF.WithFields(log.Fields{"length": len(protoMts.Metrics)}).Debug("Metrics chunk received from Snap")

		streamMts := convertProtoToMetrics(protoMts.Metrics)
		mts = append(mts, streamMts...)
		cfg = fromProtoConfig(protoMts.Config)
	}

	errMsg := ""
	if len(mts) != 0 {
		logF.WithFields(log.Fields{"length": len(mts)}).Debug("Metrics will be published")

		err := p.plugin.Publish(mts, cfg)

		if err != nil {
			errMsg = err.Error()
		}
	} else {
		logF.Info("Nothing to publish, request will be ignored.")
	}

	reply := &rpc.ErrReply{Error: errMsg}
	return stream.SendAndClose(reply)
}
