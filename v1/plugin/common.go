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
	"github.com/solarwinds/snap-plugin-lib/v1/plugin/rpc"
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
