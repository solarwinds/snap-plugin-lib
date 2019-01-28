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
	"fmt"
	"io"
	"strings"

	"github.com/librato/snap-plugin-lib-go/v1/plugin/rpc"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

//TODO(danielscottt): plugin panics

type publisherProxy struct {
	pluginProxy

	plugin Publisher
}

func (p *publisherProxy) Publish(ctx context.Context, arg *rpc.PubProcArg) (*rpc.ErrReply, error) {
	var logF = logrus.WithFields(logrus.Fields{"function": "Publish", "layer": "lib-go"})

	mts := convertProtoToMetrics(arg.Metrics)
	cfg := fromProtoConfig(arg.Config)

	logF.Debugf("Metrics will be sent to appoptics (len=%d)", len(mts))
	err := p.plugin.Publish(mts, cfg)
	if err != nil {
		return &rpc.ErrReply{Error: err.Error()}, nil
	}
	return &rpc.ErrReply{}, nil
}

func (p *publisherProxy) PublishAsStream(stream rpc.Publisher_PublishAsStreamServer) error {
	var logF = logrus.WithFields(logrus.Fields{"function": "PublishAsStream", "layer": "lib-go"})

	var errList []string

	for {
		protoMts, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			return fmt.Errorf("failure when reading from stream: %s", err.Error())
		}

		logF.Debugf("Metrics chunk will be sent to appoptics (len=%d)", len(protoMts.Metrics))

		mts := convertProtoToMetrics(protoMts.Metrics)
		cfg := fromProtoConfig(protoMts.Config)

		err = p.plugin.Publish(mts, cfg)
		if err != nil {
			errList = append(errList, err.Error())
		}
	}

	reply := &rpc.ErrReply{Error: strings.Join(errList, "")}
	return stream.SendAndClose(reply)
}
