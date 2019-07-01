package runner

import (
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

func startPprofServer(opt *plugin.Options, ln net.Listener) {
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

func startStatsServer(_ *plugin.Options) {
	log.Infof("Running stats server on address")

	// TODO: AO-13450
}
