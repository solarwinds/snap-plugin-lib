/*
Package rpc:
* contains Protocol Buffer types definitions
* handles GRPC communication (server side), passing it to proxies.
* contains Implementation of GRPC services.
*/
package pluginrpc

import (
	"net"
	"time"

	"github.com/fullstorydev/grpchan"

	"github.com/librato/snap-plugin-lib-go/v2/internal/plugins/common/stats"
	"github.com/librato/snap-plugin-lib-go/v2/pluginrpc"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const GRPCGracefulStopTimeout = 10 * time.Second

var log = logrus.WithFields(logrus.Fields{"layer": "lib", "module": "plugin-rpc"})

type Server interface {
	grpchan.ServiceRegistry

	// For compatibility with the native grpc.Server
	Serve(lis net.Listener) error
	GracefulStop()
	Stop()
}

func NewGRPCServer(inProc bool) Server {
	if inProc {
		return NewChannel()
	}

	return grpc.NewServer()
}

func StartCollectorGRPC(srv Server, proxy CollectorProxy, statsController stats.Controller, grpcLn net.Listener, pprofLn net.Listener, pingTimeout time.Duration, pingMaxMissedCount uint) {
	pluginrpc.RegisterHandlerCollector(srv, newCollectService(proxy, statsController, pprofLn))
	startGRPC(srv, grpcLn, pingTimeout, pingMaxMissedCount)
}

func StartPublisherGRPC(srv Server, proxy PublisherProxy, statsController stats.Controller, grpcLn net.Listener, pprofLn net.Listener, pingTimeout time.Duration, pingMaxMissedCount uint) {
	pluginrpc.RegisterHandlerPublisher(srv, newPublishingService(proxy, statsController, pprofLn))
	startGRPC(srv, grpcLn, pingTimeout, pingMaxMissedCount)
}

func startGRPC(srv Server, grpcLn net.Listener, pingTimeout time.Duration, pingMaxMissedCount uint) {
	closeChan := make(chan error, 1)
	pluginrpc.RegisterHandlerController(srv, newControlService(closeChan, pingTimeout, pingMaxMissedCount))

	go func() {
		err := srv.Serve(grpcLn) // may be blocking (depending on implementation)
		if err != nil {
			closeChan <- err
		}
	}()

	exitErr := <-closeChan
	if exitErr != nil && exitErr != RequestedKillError {
		log.WithError(exitErr).Errorf("Major error occurred - plugin will be shut down")
	}

	shutdownPlugin(srv)
}

func shutdownPlugin(srv Server) {
	stopped := make(chan bool, 1)

	// try to complete all remaining rpc calls
	go func() {
		srv.GracefulStop()
		stopped <- true
	}()

	// If RPC calls lasting too much, stop server by force
	select {
	case <-stopped:
		log.Debug("GRPC server stopped gracefully")
	case <-time.After(GRPCGracefulStopTimeout):
		srv.Stop()
		log.Warning("GRPC server couldn't have been stopped gracefully. Some metrics might be lost")
	}
}
