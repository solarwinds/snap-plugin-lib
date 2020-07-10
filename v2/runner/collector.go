/*
The package "runner" provides simple API to start plugins in different modes.
*/

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
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/collector/proxy"
	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/stats"
	"github.com/librato/snap-plugin-lib-go/v2/internal/service"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/log"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/types"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/sirupsen/logrus"
)

const (
	normalExitStatus = 0
	errorExitStatus  = 1

	infiniteDebugCollectCount = -1
)

func StartStreamingCollector(collector plugin.StreamingCollector, name string, version string) {
	StartStreamingCollectorWithContext(context.Background(), collector, name, version)
}

func StartStreamingCollectorWithContext(ctx context.Context, collector plugin.StreamingCollector, name string, version string) {
	startCollector(ctx, types.NewStreamingCollector(name, version, collector))
}

func StartCollector(collector plugin.Collector, name string, version string) {
	StartCollectorWithContext(context.Background(), collector, name, version)
}

func StartCollectorWithContext(ctx context.Context, collector plugin.Collector, name string, version string) {
	startCollector(ctx, types.NewCollector(name, version, collector))
}

func startCollector(ctx context.Context, collector types.Collector) {
	var err error

	var opt *plugin.Options
	inprocPlugin, inProc := collector.Unwrap().(inProcessPlugin)
	if inProc {
		opt = inprocPlugin.Options()

		logger := inprocPlugin.Logger()
		ctx = log.ToCtx(ctx, logger)
	}

	logF := logger(ctx).WithField("service", "collector")

	if opt == nil {
		opt, err = ParseCmdLineOptions(os.Args[0], types.PluginTypeCollector, os.Args[1:])
		if err != nil {
			logF.WithError(err).Error("Error occured during plugin startup")
			os.Exit(errorExitStatus)
		}
	}

	err = ValidateOptions(opt)
	if err != nil {
		logF.WithError(err).Error("Invalid plugin options")
		os.Exit(errorExitStatus)
	}

	statsController, err := stats.NewController(ctx, collector.Name(), collector.Version(), collector.Type(), opt)
	if err != nil {
		logF.WithError(err).Error("Error occured when starting statistics controller")
		os.Exit(errorExitStatus)
	}

	ctxMan := proxy.NewContextManager(ctx, collector, statsController)

	logrus.SetLevel(opt.LogLevel)

	if opt.PrintVersion {
		printVersion(collector.Name(), collector.Version())
		os.Exit(normalExitStatus)
	}

	if opt.PrintExampleTask {
		printExampleTask(ctxMan, collector.Name())
		os.Exit(normalExitStatus)
	}

	if opt.DebugMode {
		startCollectorInDebugMode(ctxMan, opt)
	} else {
		r, err := acquireResources(opt)
		if err != nil {
			logF.WithError(err).Error("Can't acquire resources for plugin services")
			os.Exit(errorExitStatus)
		}

		jsonMeta := metaInformation(ctx, collector.Name(), collector.Version(), collector.Type(), opt, r, ctxMan.TasksLimit, ctxMan.InstancesLimit)
		if inProc {
			inprocPlugin.MetaChannel() <- jsonMeta
			close(inprocPlugin.MetaChannel())
		}

		if opt.EnableProfiling {
			startPprofServer(ctx, r.pprofListener)
			defer r.pprofListener.Close() // close pprof service when GRPC service has been shut down
		}

		if opt.EnableStatsServer {
			startStatsServer(ctx, r.statsListener, statsController)
			defer r.statsListener.Close() // close stats service when GRPC service has been shut down
		}

		srv, err := service.NewGRPCServer(ctx, opt)
		if err != nil {
			logF.WithError(err).Error("Can't initialize GRPC Server")
			os.Exit(errorExitStatus)
		}

		// We need to bind the gRPC client on the other end to the same channel so need to return it from here
		if inProc {
			inprocPlugin.GRPCChannel() <- srv.(*service.Channel).Channel
		}

		// main blocking operation
		service.StartCollectorGRPC(ctx, srv, ctxMan, r.grpcListener, opt.GRPCPingTimeout, opt.GRPCPingMaxMissed)
	}
}

func startCollectorInDebugMode(ctxManager *proxy.ContextManager, opt *plugin.Options) {
	const debugModeTaskID = "task-1"

	// Load task based on command line options
	var filter []string
	if opt.PluginFilter != defaultFilter {
		filter = strings.Split(opt.PluginFilter, filterSeparator)
	}

	errLoad := ctxManager.LoadTask(debugModeTaskID, []byte(opt.PluginConfig), filter)
	if errLoad != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Couldn't load a task in a standalone mode (reason: %v)\n", errLoad)
		os.Exit(errorExitStatus)
	}

	for runCount := 0; ; {
		// Request metrics collection
		chunkCh := ctxManager.RequestCollect(debugModeTaskID)

		for chunk := range chunkCh {
			if chunk.Err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error occurred during metrics collection in a standalone mode (reason: %v)\n", chunk.Err)
				os.Exit(errorExitStatus)
			}

			// Print out metrics
			fmt.Printf("Gathered metrics (length=%d): \n", len(chunk.Metrics))
			for _, mt := range chunk.Metrics {
				fmt.Printf("%s\n", mt)
			}
			fmt.Printf("\n")
		}

		// wait to request new collection or exit
		if opt.DebugCollectCounts != infiniteDebugCollectCount {
			runCount++
			if runCount == opt.DebugCollectCounts {
				break
			}
		}

		time.Sleep(opt.DebugCollectInterval)
	}

	errUnload := ctxManager.UnloadTask(debugModeTaskID)
	if errUnload != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Couldn't unload a task in a standalone mode (reason: %v)\n", errUnload)
		os.Exit(errorExitStatus)
	}
}
