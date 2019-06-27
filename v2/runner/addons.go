package runner

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
)

func startPprofServer(opt *plugin.Options) {
	log.Infof("Running profiling server on address %s:%d", opt.GrpcIp, opt.PprofPort)

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
		http.ListenAndServe(fmt.Sprintf("%s:%d", opt.GrpcIp, opt.PprofPort), h)
	}()
}
