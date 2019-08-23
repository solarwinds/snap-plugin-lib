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

	"github.com/librato/snap-plugin-lib-go/v2/internal/stats"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const GRPCGracefulStopTimeout = 10 * time.Second

var log = logrus.WithFields(logrus.Fields{"layer": "lib", "module": "plugin-rpc"})

func StartGRPCController(proxy CollectorProxy, statsController stats.Controller, grpcLn net.Listener, pprofLn net.Listener, pingTimeout time.Duration, pingMaxMissedCount uint) {
	closeChan := make(chan error, 1)

	grpcServer := grpc.NewServer()
	RegisterControllerServer(grpcServer, newControlService(closeChan, pingTimeout, pingMaxMissedCount))
	RegisterCollectorServer(grpcServer, newCollectService(proxy, statsController, pprofLn))

	go func() {
		err := grpcServer.Serve(grpcLn)
		if err != nil {
			closeChan <- err
		}
	}()

	exitErr := <-closeChan
	if exitErr != nil && exitErr != RequestedKillError {
		log.WithError(exitErr).Errorf("Major error occurred - plugin will be shut down")
	}

	shutdownPlugin(grpcServer)
}

func shutdownPlugin(grpcServer *grpc.Server) {
	stopped := make(chan bool, 1)

	// try to complete all remaining rpc calls
	go func() {
		grpcServer.GracefulStop()
		stopped <- true
	}()

	// If RPC calls lasting too much, stop server by force
	select {
	case <-stopped:
		log.Debug("GRPC server stopped gracefully")
	case <-time.After(GRPCGracefulStopTimeout):
		grpcServer.Stop()
		log.Warning("GRPC server couldn't have been stopped gracefully. Some metrics might be lost")
	}
}
