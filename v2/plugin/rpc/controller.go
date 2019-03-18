package rpc

import (
	"github.com/librato/snap-plugin-lib-go/v2/plugin/proxy"
	"google.golang.org/grpc"
	"net"
)

type GRPCController struct {
	CloseChan chan error
}

func StartGRPCController(proxy proxy.Collector) {
	closeChan := make(chan error, 1)

	lis, err := net.Listen("tcp", "0.0.0.0:56789")
	if err != nil {
		// todo: raise or log an error
	}

	grpcServer := grpc.NewServer()
	RegisterControllerServer(grpcServer, newControlService(closeChan))
	RegisterCollectorServer(grpcServer, newCollectService(proxy))

	go func() {
		err := grpcServer.Serve(lis)
		if err != nil {
			closeChan <- err
		}
	}()

	<-closeChan
}
