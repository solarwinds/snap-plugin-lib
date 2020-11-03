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
	"os"

	"github.com/sirupsen/logrus"

	"github.com/solarwinds/snap-plugin-lib/v2/internal/plugins/common/stats"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/plugins/publisher/proxy"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/service"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/log"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/types"
	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
)

func StartPublisher(publisher plugin.Publisher, name string, version string) {
	StartPublisherWithContext(context.Background(), publisher, name, version)
}

func StartPublisherWithContext(ctx context.Context, publisher plugin.Publisher, name string, version string) {
	var err error
	var opt *plugin.Options
	var ctxLog logrus.FieldLogger

	inprocPlugin, inProc := publisher.(inProcessPlugin)
	if inProc {
		opt = inprocPlugin.Options()
		ctxLog = inprocPlugin.Logger()
		publisher = inprocPlugin.Unwrap().(plugin.Publisher)
	} else {
		ctxLog = logrus.WithFields(logrus.Fields{
			"plugin-name":    name,
			"plugin-version": version,
		})
	}

	ctx = log.ToCtx(ctx, ctxLog)

	logF := logger(ctx).WithField("service", "publisher")

	if opt == nil {
		opt, err = ParseCmdLineOptions(os.Args[0], types.PluginTypePublisher, os.Args[1:])
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

	statsController, err := stats.NewController(ctx, name, version, types.PluginTypePublisher, opt)
	if err != nil {
		logF.WithError(err).Error("Error occured when starting statistics controller")
		os.Exit(errorExitStatus)
	}

	ctxMan := proxy.NewContextManager(publisher, statsController)

	logrus.SetLevel(opt.LogLevel)

	if opt.PrintVersion {
		printVersion(name, version)
		os.Exit(normalExitStatus)
	}

	r, err := acquireResources(opt)
	if err != nil {
		logF.WithError(err).Error("Can't acquire resources for plugin services")
		os.Exit(errorExitStatus)
	}

	jsonMeta := metaInformation(ctx, name, version, types.PluginTypePublisher, opt, r, ctxMan.TasksLimit, ctxMan.InstancesLimit)
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
	service.StartPublisherGRPC(ctx, srv, ctxMan, r.grpcListener, opt.GRPCPingTimeout, opt.GRPCPingMaxMissed)
}
