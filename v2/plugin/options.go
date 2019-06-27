package plugin

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Options struct {
	GrpcIp            string
	GrpcPort          int
	GrpcPingTimeout   time.Duration
	GrpcPingMaxMissed int

	LogLevel     logrus.Level
	EnablePprof  bool
	PprofPort    int
	EnableStats  bool
	DebugMode    bool
	PluginConfig string
}
