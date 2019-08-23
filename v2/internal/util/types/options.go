package types

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Structure representing plugin configuration (received by parsing command-line arguments)
// Visit newFlagParser() to find descriptions associated with each option.
type Options struct {
	PluginIp          string
	GrpcPort          int
	GrpcPingTimeout   time.Duration
	GrpcPingMaxMissed uint

	LogLevel          logrus.Level
	EnableProfiling   bool
	EnableStats       bool // enable calculation statistics
	EnableStatsServer bool // if true, start statistics HTTP server
	PprofPort         int  `json:",omitempty"`
	StatsPort         int  `json:",omitempty"`

	PrintExampleTask     bool          `json:"-"`
	DebugMode            bool          `json:"-"`
	PluginConfig         string        `json:"-"`
	PluginFilter         string        `json:"-"`
	DebugCollectCounts   uint          `json:"-"`
	DebugCollectInterval time.Duration `json:"-"`
}
