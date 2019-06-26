package plugin

import "github.com/sirupsen/logrus"

type Options struct {
	GrpcIp   string
	GrpcPort int

	LogLevel     logrus.Level
	EnablePprof  bool
	EnableStats  bool
	DebugMode    bool
	PluginConfig string
}
