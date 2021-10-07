/*
 Copyright (c) 2021 SolarWinds Worldwide, LLC

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
	"time"

	"github.com/sirupsen/logrus"
)

// Structure representing plugin configuration (received by parsing command-line arguments)
// Visit newFlagParser() to find descriptions associated with each option.
type Options struct {
	PluginIP          string
	GRPCPort          int
	GRPCPingTimeout   time.Duration
	GRPCPingMaxMissed uint

	EnableTLS         bool // GRPC Server
	TLSServerCertPath string
	TLSServerKeyPath  string
	TLSClientCAPath   string

	LogLevel          logrus.Level
	EnableProfiling   bool
	PProfPort         int  `json:",omitempty"`
	EnableStats       bool // enable calculation statistics
	EnableStatsServer bool // if true, start statistics HTTP server
	StatsPort         int  `json:",omitempty"`

	UseAPIv2 bool
	AsThread bool

	CollectChunkSize uint

	PrintExampleTask     bool          `json:"-"`
	DebugMode            bool          `json:"-"`
	PluginConfig         string        `json:"-"`
	PluginFilter         string        `json:"-"`
	DebugCollectCounts   int           `json:"-"`
	DebugCollectInterval time.Duration `json:"-"`
	PrintVersion         bool          `json:"-"`
}
