/*
Package rpc:
* contains Protocol Buffer types definitions
* handles GRPC communication (server side), passing it to proxies.
* contains Implementation of GRPC services.
*/
package service

import (
	"context"
	"net"
	"time"

	"github.com/librato/grpchan"
	"github.com/librato/snap-plugin-lib-go/v2/internal/util/log"
	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/pluginrpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const GRPCGracefulStopTimeout = 10 * time.Second

var moduleFields = logrus.Fields{"layer": "lib", "module": "plugin-rpc"}

type Server interface {
	grpchan.ServiceRegistry

	// For compatibility with the native grpc.Server
	Serve(lis net.Listener) error
	GracefulStop()
	Stop()
}

// An abstraction providing a unified interface for
// * the native go-grpc implementation
// * https://github.com/librato/grpchan - this one provides a way of using gRPC with a custom transport
//   (that means sth other than the native h2 - HTTP1.1 or inprocess/channels are available out of the box)
func NewGRPCServer(ctx context.Context, opt *plugin.Options) (Server, error) {
	if opt.AsThread {
		return NewChannel(), nil
	}

	if !opt.EnableTLS {
		return grpc.NewServer(), nil
	}

	tlsCreds, err := tlsCredentials(ctx, opt)
	if err != nil {
		return nil, err
	}

	return grpc.NewServer(grpc.Creds(tlsCreds)), nil
}

func StartCollectorGRPC(ctx context.Context, srv Server, proxy CollectorProxy, grpcLn net.Listener, pingTimeout time.Duration, pingMaxMissedCount uint) {
	pluginrpc.RegisterHandlerCollector(srv, newCollectService(ctx, proxy))
	startGRPC(ctx, srv, grpcLn, pingTimeout, pingMaxMissedCount)
}

func StartPublisherGRPC(ctx context.Context, srv Server, proxy PublisherProxy, grpcLn net.Listener, pingTimeout time.Duration, pingMaxMissedCount uint) {
	pluginrpc.RegisterHandlerPublisher(srv, newPublishingService(ctx, proxy))
	startGRPC(ctx, srv, grpcLn, pingTimeout, pingMaxMissedCount)
}

func startGRPC(ctx context.Context, srv Server, grpcLn net.Listener, pingTimeout time.Duration, pingMaxMissedCount uint) {
	errChan := make(chan error)

	csCtx, cancelFn := context.WithCancel(ctx)
	pluginrpc.RegisterHandlerController(srv, newControlService(csCtx, errChan, pingTimeout, pingMaxMissedCount))

	go func() {
		err := srv.Serve(grpcLn) // may be blocking (depending on implementation)
		if err != nil {
			errChan <- err
		}
	}()

	err := <-errChan // may be blocking (depending on implementation)
	cancelFn()       // signal ping monitor (via ctx)

	if err != nil && err != RequestedKillError {
		log.WithCtx(ctx).WithFields(moduleFields).
			WithError(err).Errorf("Major error occurred - plugin will be shut down")
	}

	shutdownPlugin(ctx, srv)
}

func shutdownPlugin(ctx context.Context, srv Server) {
	stopped := make(chan bool, 1)

	// try to complete all remaining rpc calls
	go func() {
		srv.GracefulStop()
		stopped <- true
	}()

	logF := log.WithCtx(ctx).WithFields(moduleFields)

	// If RPC calls lasting too much, stop server by force
	select {
	case <-stopped:
		logF.Debug("GRPC server stopped gracefully")
	case <-time.After(GRPCGracefulStopTimeout):
		srv.Stop()
		logF.Warning("GRPC server couldn't have been stopped gracefully. Some metrics might have been lost")
	}
}
