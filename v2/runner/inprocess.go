/*
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

package runner

import (
	"github.com/sirupsen/logrus"
	"github.com/solarwinds/grpchan"

	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
)

type inProcessPlugin interface {
	plugin.InProcessPlugin
	Unwrap() interface{}

	Options() *plugin.Options
	GRPCChannel() chan<- grpchan.Channel
	MetaChannel() chan<- []byte
	Logger() logrus.FieldLogger
}
