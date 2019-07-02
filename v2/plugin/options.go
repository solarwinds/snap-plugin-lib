package plugin

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Options struct {
	PluginIp          string
	GrpcPort          int
	GrpcPingTimeout   time.Duration
	GrpcPingMaxMissed uint

	LogLevel    logrus.Level
	EnablePprof bool
	EnableStats bool
	PprofPort   int
	StatsPort   int

	DebugMode            bool
	PluginConfig         string
	PluginFilter         string
	DebugCollectCounts   uint
	DebugCollectInterval time.Duration
}
