// +build small

package service

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/plugin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConnectingToSecureGRPC(t *testing.T) {
	Convey("Validate that client can connect to secure GRPC Server", t, func() {
		// Arrange
		opt := &plugin.Options{
			EnableTLS:   true,
			TLSKeyPath:  filepath.Join("cert_test", "serv.key"),
			TLSCertPath: filepath.Join("cert_test", "serv.crt"),
		}

		var ln net.Listener

		// Arrange (GRPC Server)
		go func() {
			srv, _ := NewGRPCServer(opt)

			ln, _ = net.Listen("tcp", "")
			_ = srv.Serve(ln)
		}()

		time.Sleep(2 * time.Second)

		// Arrange (GRPC Client)
		caCertPath := filepath.Join("cert_test", "ca.crt")
		caCert, _ := ioutil.ReadFile(caCertPath)
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(caCert)

		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: false,
			RootCAs: certPool,
		})
		conn, err := grpc.Dial(ln.Addr().String(), grpc.WithTransportCredentials(creds))
		fmt.Printf("conn %v\n", conn)
		fmt.Printf("err %v\n", err)

		time.Sleep(10 * time.Second)

	})
}
