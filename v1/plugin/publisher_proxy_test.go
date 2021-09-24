//go:build small
// +build small

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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/solarwinds/snap-plugin-lib/v1/plugin/rpc"
)

func TestPublisher(t *testing.T) {
	Convey("Test Publisher", t, func() {
		Convey("Error while publishing", func() {
			pp := publisherProxy{
				pluginProxy: *newPluginProxy(newMockPublisher()),
				plugin:      newMockErrPublisher(),
			}
			errReply, err := pp.Publish(context.Background(), &rpc.PubProcArg{})
			So(errReply.Error, ShouldResemble, "error")
			So(err, ShouldBeNil)
		})
		Convey("Succeed while publishing", func() {
			pp := publisherProxy{
				pluginProxy: *newPluginProxy(newMockPublisher()),
				plugin:      newMockPublisher(),
			}

			input, err := getTestData()
			So(err, ShouldBeNil)

			_, err = pp.Publish(context.Background(), &rpc.PubProcArg{Metrics: input})
			So(err, ShouldBeNil)
		})
	})
}

func getTestData() ([]*rpc.Metric, error) {
	input := []*rpc.Metric{}

	mp := getMockMetricDataMap()
	for _, v := range mp {
		m, err := toProtoMetric(v)
		if err != nil {
			return nil, err
		}
		input = append(input, m)
	}
	return input, nil
}
