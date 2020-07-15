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
	"flag"
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/service"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/sirupsen/logrus"
)

const (
	defaultPluginIP  = "127.0.0.1"
	defaultGRPCPort  = 0
	defaultPProfPort = 0
	defaultStatsPort = 0

	defaultConfig          = "{}"
	defaultFilter          = ""
	defaultCollectInterval = 5 * time.Second
	defaultCollectCount    = 1

	defaultLogLevel = logrus.WarnLevel

	filterSeparator = ";"
)

///////////////////////////////////////////////////////////////////////////////

func newFlagParser(name string, pType types.PluginType, opt *plugin.Options) *flag.FlagSet {
	flagParser := flag.NewFlagSet(name, flag.ContinueOnError)

	// common flags
	flagParser.BoolVar(&opt.PrintVersion,
		"version", false,
		"Print version of plugin")

	flagParser.StringVar(&opt.PluginIP,
		"plugin-ip", defaultPluginIP,
		"IP Address on which GRPC server will be served")

	flagParser.IntVar(&opt.GRPCPort,
		"grpc-port", defaultGRPCPort,
		"Port on which GRPC server will be served")

	flagParser.DurationVar(&opt.GRPCPingTimeout,
		"grpc-ping-timeout", service.DefaultPingTimeout,
		"Deadline for receiving single ping messages")

	flagParser.UintVar(&opt.GRPCPingMaxMissed,
		"grpc-ping-max-missed", service.DefaultMaxMissingPingCounter,
		"Number of missed ping messages after which plugin should exit")

	allLogLevels := strings.Replace(fmt.Sprintf("%v", logrus.AllLevels), " ", ", ", -1)
	flagParser.Var(&logLevelHandler{opt: opt},
		"log-level",
		fmt.Sprintf("Minimal level of logged messages %s", allLogLevels))

	flagParser.BoolVar(&opt.EnableProfiling,
		"enable-profiling", false,
		"Enable profiling (pprof server)")

	flagParser.IntVar(&opt.PProfPort,
		"pprof-port", defaultPProfPort,
		"Port on which profiling server will be available")

	flagParser.BoolVar(&opt.EnableStats,
		"enable-stats", false,
		"Enable gathering plugin statistics")

	flagParser.BoolVar(&opt.EnableStatsServer,
		"enable-stats-server", false,
		"Enable stats server")

	flagParser.IntVar(&opt.StatsPort,
		"stats-port", defaultStatsPort,
		"Port on which stats server will be available")

	flagParser.BoolVar(&opt.UseAPIv2,
		"plugin-api-v2", true,
		"If a plugin supports multiple plugin API versions, set it to use v2")

	flagParser.BoolVar(&opt.EnableTLS,
		"tls", false,
		"Enable secure GRPC communication")

	flagParser.StringVar(&opt.TLSServerCertPath,
		"cert-path", "",
		"Certificate path used by GRPC Server")

	flagParser.StringVar(&opt.TLSServerKeyPath,
		"key-path", "",
		"Path to private key associated with server certificate")

	flagParser.StringVar(&opt.TLSClientCAPath,
		"root-cert-paths", "",
		fmt.Sprintf("Path to CA root path certificate(s). Might also be provided as files or/and dirs separated with '%c'.", filepath.Separator))

	// custom flags

	if pType == types.PluginTypeCollector {
		flagParser.BoolVar(&opt.PrintExampleTask,
			"print-example-task", false,
			"Print-out example task for a plugin")

		flagParser.BoolVar(&opt.DebugMode,
			"debug-mode", false,
			"Run plugin in debug mode (standalone)")

		flagParser.IntVar(&opt.DebugCollectCounts,
			"debug-collect-counts", defaultCollectCount,
			"Number of collect requests executed in debug mode (-1 for infinitely)")

		flagParser.DurationVar(&opt.DebugCollectInterval,
			"debug-collect-interval", defaultCollectInterval,
			"Interval between consecutive collect requests")

		flagParser.StringVar(&opt.PluginConfig,
			"plugin-config", defaultConfig,
			"Collector configuration in debug mode")

		flagParser.StringVar(&opt.PluginFilter,
			"plugin-filter", defaultFilter,
			fmt.Sprintf("Default filtering definition (separated by %s)", filterSeparator))
	}

	return flagParser
}

type logLevelHandler struct {
	opt *plugin.Options
}

func (l *logLevelHandler) String() string {
	if l.opt == nil {
		return "error"
	}

	return l.opt.LogLevel.String()
}

func (l *logLevelHandler) Set(s string) error {
	// accept level as a form of int (0 - 6)
	intLvl, errConv := strconv.Atoi(s)
	if errConv == nil && intLvl >= int(logrus.PanicLevel) && intLvl <= int(logrus.TraceLevel) {
		l.opt.LogLevel = logrus.Level(intLvl)
		return nil
	}

	// accept level as a form os string (warning, error etc.)
	lvl, errParse := logrus.ParseLevel(s)
	if errParse != nil {
		return errParse
	}
	l.opt.LogLevel = lvl

	return nil
}

///////////////////////////////////////////////////////////////////////////////

func ParseCmdLineOptions(pluginName string, pluginType types.PluginType, args []string) (*plugin.Options, error) {
	opt := &plugin.Options{
		LogLevel: defaultLogLevel,
	}

	flagParser := newFlagParser(pluginName, pluginType, opt)
	argsToParse := args[:]
	if len(args) > 0 && strings.HasSuffix(args[0], ".py") {
		// ignore first parameter if plugin is an interpreted code
		argsToParse = args[1:]
	}

	err := flagParser.Parse(argsToParse)
	if err != nil {
		return opt, fmt.Errorf("can't parse command line options: %v", err)
	}

	v := flagParser.Args()
	if len(v) > 0 {
		return opt, fmt.Errorf("unexpected option(s) provided: %v %v", v, len(v))
	}

	return opt, nil
}

func ValidateOptions(opt *plugin.Options) error {
	if opt.DebugCollectCounts == 0 {
		opt.DebugCollectCounts = defaultCollectCount
	}

	if opt.DebugCollectInterval == 0 {
		opt.DebugCollectInterval = defaultCollectInterval
	}

	if opt.PluginConfig == "" {
		opt.PluginConfig = defaultConfig
	}

	if opt.PluginFilter == "" {
		opt.PluginFilter = defaultFilter
	}

	grpcIp := net.ParseIP(opt.PluginIP)
	if grpcIp == nil {
		return fmt.Errorf("GRPC IP contains invalid address")
	}

	if opt.EnableTLS {
		if opt.TLSServerCertPath == "" || opt.TLSServerKeyPath == "" {
			return fmt.Errorf("certificate and key path have to be provided when TLS is enabled")
		}
	}

	if opt.PProfPort > 0 && !opt.EnableProfiling {
		return fmt.Errorf("-enable-pprof flag should be set when configuring pprof port")
	}

	if opt.StatsPort > 0 && !opt.EnableStatsServer {
		return fmt.Errorf("-enable-stats flag should be set when configuring stats port")
	}

	if opt.EnableStatsServer && !opt.EnableStats {
		return fmt.Errorf("-enable-stats should be set when -enable-stats-server=1")
	}

	if !opt.DebugMode && anyDebugFlagSet(opt) {
		return fmt.Errorf("-debug-mode flag should be set when configuring debug options")
	}

	return nil
}

func anyDebugFlagSet(opt *plugin.Options) bool {
	return opt.DebugCollectCounts != defaultCollectCount ||
		opt.DebugCollectInterval != defaultCollectInterval ||
		opt.PluginConfig != defaultConfig ||
		opt.PluginFilter != defaultFilter
}
