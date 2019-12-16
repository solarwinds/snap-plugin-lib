package runner

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internals/plugins/common/stats"
)

const (
	statsRequestTimeout = 10 * time.Second

	jsonIndentString = "    "
)

///////////////////////////////////////////////////////////////////////////////

func startPprofServer(ln net.Listener) {
	log.Infof("Running profiling server on address %s", ln.Addr())

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
			log.WithError(err).Warn("Pprof server stopped")
		}
	}()
}

///////////////////////////////////////////////////////////////////////////////

func startStatsServer(ln net.Listener, stats stats.Controller) {
	log.Infof("Running stats server on address")

	h := http.NewServeMux()

	h.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		statsHandler(w, r, stats)
	})

	go func() {
		err := http.Serve(ln, h)
		if err != nil {
			log.WithError(err).Warn("Stats server stopped")
		}
	}()
}

func statsHandler(w http.ResponseWriter, r *http.Request, stats stats.Controller) {
	log.WithField("URI", r.RequestURI).Trace("Handling statistics request")

	respCh := stats.RequestStat()

	select {
	case resp := <-respCh:
		jsonStats, err := json.MarshalIndent(resp, "", jsonIndentString)
		if err != nil {
			log.WithField("stats", fmt.Sprintf("%v", resp)).WithError(err).Error("error when marshaling statistics struct")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonStats)
		if err != nil {
			log.WithError(err).Error("error occurred when serving statistics request")
		}

	case <-time.After(statsRequestTimeout):
		log.WithField("timeout", statsRequestTimeout).Warning("timeout occurred when serving statistics request")
		w.WriteHeader(http.StatusRequestTimeout)
	}
}
