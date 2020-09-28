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
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/solarwinds/snap-plugin-lib/v2/internal/plugins/common/stats"
	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/log"
)

const (
	statsRequestTimeout = 10 * time.Second

	jsonIndentString = "    "
)

///////////////////////////////////////////////////////////////////////////////

func startPprofServer(ctx context.Context, ln net.Listener) {
	logF := log.WithCtx(ctx).WithFields(moduleFields)
	logF.Infof("Running profiling server on address %s", ln.Addr())

	h := http.NewServeMux()

	h.HandleFunc("/debug/pprof/", func(w http.ResponseWriter, r *http.Request) {
		pprof.Index(w, r)
	})
	h.HandleFunc("/debug/pprof/profile", func(w http.ResponseWriter, r *http.Request) {
		pprof.Profile(w, r)
	})
	h.HandleFunc("/debug/pprof/trace", func(w http.ResponseWriter, r *http.Request) {
		pprof.Trace(w, r)
	})
	h.HandleFunc("/debug/pprof/cmdline", func(w http.ResponseWriter, r *http.Request) {
		pprof.Cmdline(w, r)
	})

	go func() {
		err := http.Serve(ln, h)
		if err != nil {
			logF.WithError(err).Warn("Pprof server stopped")
		}
	}()
}

///////////////////////////////////////////////////////////////////////////////

func startStatsServer(ctx context.Context, ln net.Listener, stats stats.Controller) {
	logF := log.WithCtx(ctx).WithFields(moduleFields)
	logF.Infof("Running stats server on address %s", ln.Addr())

	h := http.NewServeMux()

	h.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		statsHandler(ctx, w, r, stats)
	})

	go func() {
		err := http.Serve(ln, h)
		if err != nil {
			logF.WithError(err).Warn("Stats server stopped")
		}
	}()
}

func statsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, stats stats.Controller) {
	logF := log.WithCtx(ctx).WithFields(moduleFields)
	logF.WithField("URI", r.RequestURI).Trace("Handling statistics request")

	respCh := stats.RequestStat()

	select {
	case resp := <-respCh:
		jsonStats, err := json.MarshalIndent(resp, "", jsonIndentString)
		if err != nil {
			logF.WithField("stats", fmt.Sprintf("%v", resp)).WithError(err).Error("error when marshaling statistics struct")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonStats)
		if err != nil {
			logF.WithError(err).Error("error occurred when serving statistics request")
		}

	case <-time.After(statsRequestTimeout):
		logF.WithField("timeout", statsRequestTimeout).Warn("timeout occurred when serving statistics request")
		w.WriteHeader(http.StatusRequestTimeout)
	}
}
