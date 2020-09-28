// +build medium

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

package service

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
	"github.com/solarwinds/snap-plugin-lib/v2/pluginrpc"
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
			srv, _ := NewGRPCServer(context.Background(), opt)
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
