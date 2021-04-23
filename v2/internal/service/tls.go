/*
 Copyright (c) 2021 SolarWinds Worldwide, LLC

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
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/log"
	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
	"google.golang.org/grpc/credentials"
)

func tlsCredentials(ctx context.Context, opt *plugin.Options) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(opt.TLSServerCertPath, opt.TLSServerKeyPath)
	if err != nil {
		return nil, fmt.Errorf("invalid TLS certificate: %v", err)
	}

	clientCA, err := loadCACerts(ctx, opt.TLSClientCAPath)
	if err != nil {
		return nil, fmt.Errorf("can't read client CA Cert(s)")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCA,
		MinVersion:   tls.VersionTLS12,
	}
	tlsCreds := credentials.NewTLS(tlsConfig)

	return tlsCreds, nil
}

func loadCACerts(ctx context.Context, caPath string) (*x509.CertPool, error) {
	logF := log.WithCtx(ctx).WithFields(moduleFields)

	clientCA := x509.NewCertPool()
	certFilesToProcess := []string{}

	// gather list of files to read
	for _, p := range filepath.SplitList(caPath) {
		info, err := os.Stat(p)
		if err != nil {
			return clientCA, fmt.Errorf("path doesn't point to any file or directory (%v): %v", p, err)
		}

		if !info.IsDir() {
			certFilesToProcess = append(certFilesToProcess, p)
			continue
		}

		_ = filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				certFilesToProcess = append(certFilesToProcess, path)
			}
			return nil
		})
	}

	// read certificate files
	for _, f := range certFilesToProcess {
		caCert, err := ioutil.ReadFile(filepath.Clean(f))
		if err != nil {
			return clientCA, fmt.Errorf("can't read client CA Root certificate: %v", err)
		}
		ok := clientCA.AppendCertsFromPEM(caCert)
		if !ok {
			logF.WithField("path", f).Warn("given file is not a certificate")
		}
	}

	return clientCA, nil
}
