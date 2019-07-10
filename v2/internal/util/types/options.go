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

	LogLevel    logrus.Level
	EnablePprof bool
	EnableStats bool
	PprofPort   int `json:",omitempty"`
	StatsPort   int `json:",omitempty"`

	DebugMode            bool          `json:"-"`
	PluginConfig         string        `json:"-"`
	PluginFilter         string        `json:"-"`
	DebugCollectCounts   uint          `json:"-"`
	DebugCollectInterval time.Duration `json:"-"`
}
