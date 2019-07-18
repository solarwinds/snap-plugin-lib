package runner

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/pluginrpc"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
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

func newFlagParser(name string, opt *types.Options) *flag.FlagSet {
	flagParser := flag.NewFlagSet(name, flag.ContinueOnError)

	flagParser.StringVar(&opt.PluginIp,
		"plugin-ip", defaultPluginIP,
		"IP Address on which GRPC server will be served")

	flagParser.IntVar(&opt.GrpcPort,
		"grpc-port", defaultGRPCPort,
		"Port on which GRPC server will be served")

	flagParser.DurationVar(&opt.GrpcPingTimeout,
		"grpc-ping-timeout", pluginrpc.DefaultPingTimeout,
		"Deadline for receiving single ping messages")

	flagParser.UintVar(&opt.GrpcPingMaxMissed,
		"grpc-ping-max-missed", pluginrpc.DefaultMaxMissingPingCounter,
		"Number of missed ping messages after which plugin should exit")

	allLogLevels := strings.Replace(fmt.Sprintf("%v", logrus.AllLevels), " ", ", ", -1)
	flagParser.Var(&logLevelHandler{opt: opt},
		"log-level",
		fmt.Sprintf("Minimal level of logged messages %s", allLogLevels))

	flagParser.BoolVar(&opt.EnablePprofServer,
		"enable-pprof-server", false,
		"Enable profiling server")

	flagParser.IntVar(&opt.PprofPort,
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

	flagParser.BoolVar(&opt.DebugMode,
		"debug-mode", false,
		"Run plugin in debug mode (standalone)")

	flagParser.UintVar(&opt.DebugCollectCounts,
		"debug-collect-counts", defaultCollectCount,
		"Number of collect requests executed in debug mode (0 - infinitely)")

	flagParser.DurationVar(&opt.DebugCollectInterval,
		"debug-collect-interval", defaultCollectInterval,
		"Interval between consecutive collect requests")

	flagParser.StringVar(&opt.PluginConfig,
		"plugin-config", defaultConfig,
		"Collector configuration in debug mode")

	flagParser.StringVar(&opt.PluginFilter,
		"plugin-filter", defaultFilter,
		fmt.Sprintf("Default filtering definition (separated by %s)", filterSeparator))

	return flagParser
}

type logLevelHandler struct {
	opt *types.Options
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

func ParseCmdLineOptions(pluginName string, args []string) (*types.Options, error) {
	opt := &types.Options{
		LogLevel: defaultLogLevel,
	}

	flagParser := newFlagParser(pluginName, opt)

	err := flagParser.Parse(args)
	if err != nil {
		return opt, fmt.Errorf("can't parse command line options: %v", err)
	}

	v := flagParser.Args()
	if len(v) > 0 {
		return opt, fmt.Errorf("unexpected option(s) provided: %v %v", v, len(v))
	}

	return opt, nil
}

func ValidateOptions(opt *types.Options) error {
	grpcIp := net.ParseIP(opt.PluginIp)
	if grpcIp == nil {
		return fmt.Errorf("GRPC IP contains invalid address")
	}

	if opt.PprofPort > 0 && !opt.EnablePprofServer {
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

func anyDebugFlagSet(opt *types.Options) bool {
	return opt.DebugCollectCounts != defaultCollectCount ||
		opt.DebugCollectInterval != defaultCollectInterval ||
		opt.PluginConfig != defaultConfig ||
		opt.PluginFilter != defaultFilter
}
