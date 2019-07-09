package runner

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/librato/snap-plugin-lib-go/v2/internal/stats"
)

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
		http.Serve(ln, h)
	}()
}

func startStatsServer(ln net.Listener, stats *stats.StatsController) {
	log.Infof("Running stats server on address")

	h := http.NewServeMux()

	h.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		log.Trace("Handling stats request")

		res := <-stats.RequestStat()
		b, err := json.MarshalIndent(&res, "", "    ")
		if err != nil {
			// todo:
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	go func() {
		http.Serve(ln, h)
	}()
}
