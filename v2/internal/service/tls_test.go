// +build medium

package service

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	"github.com/librato/snap-plugin-lib-go/v2/pluginrpc"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

///////////////////////////////////////////////////////////////////////////////

type controlMock struct {
	closeCh chan bool
}

func (c *controlMock) Ping(ctx context.Context, request *pluginrpc.PingRequest) (*pluginrpc.PingResponse, error) {
	return &pluginrpc.PingResponse{}, nil
}

func (c *controlMock) Kill(ctx context.Context, request *pluginrpc.KillRequest) (*pluginrpc.KillResponse, error) {
	close(c.closeCh)
	return &pluginrpc.KillResponse{}, nil
}

///////////////////////////////////////////////////////////////////////////////

const (
	certificateFolderName = "cert_test"
	grpcPort              = 50123
	grpcServerAddr        = "localhost:50123"

	tlsTestTimeout = 10 * time.Second
)

func TestConnectingToSecureGRPC(t *testing.T) {
	Convey("Validate that client can connect to secure GRPC Server", t, func() {
		// Arrange
		opt := &plugin.Options{
			GRPCPort:          grpcPort,
			EnableTLS:         true,
			TLSServerKeyPath:  filepath.Join(certificateFolderName, "serv.key"),
			TLSServerCertPath: filepath.Join(certificateFolderName, "serv.crt"),
			TLSClientCAPath:   filepath.Join(certificateFolderName, "ca.crt"),
		}

		ln, _ := net.Listen("tcp", grpcServerAddr)

		controlService := &controlMock{closeCh: make(chan bool)}

		// Arrange (GRPC Server)
		go func() {
			srv, _ := NewGRPCServer(opt)
			pluginrpc.RegisterControllerServer(srv.(*grpc.Server), controlService)

			go func() {
				<-controlService.closeCh
				srv.GracefulStop()
			}()

			_ = srv.Serve(ln)
		}()

		// Arrange (GRPC Client)
		caCertPath := filepath.Join(certificateFolderName, "ca.crt")
		caCert, _ := ioutil.ReadFile(caCertPath)
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(caCert)

		cliCertPath := filepath.Join(certificateFolderName, "cli.crt")
		cliKeyPath := filepath.Join(certificateFolderName, "cli.key")
		cert, _ := tls.LoadX509KeyPair(cliCertPath, cliKeyPath)

		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: false,
			RootCAs:            certPool,
			Certificates:       []tls.Certificate{cert},
		})

		// Act
		conn, dialErr := grpc.Dial(grpcServerAddr, grpc.WithTransportCredentials(creds))
		cc := pluginrpc.NewControllerClient(conn)

		_, errPing := cc.Ping(context.Background(), &pluginrpc.PingRequest{})
		_, errKill := cc.Kill(context.Background(), &pluginrpc.KillRequest{})

		// Assert
		So(dialErr, ShouldBeNil)
		So(errPing, ShouldBeNil)
		So(errKill, ShouldBeNil)

		select {
		case <-controlService.closeCh:
		// ok
		case <-time.After(tlsTestTimeout):
			panic("timeout occurred, kill wasn't properly received")
		}
	})
}
