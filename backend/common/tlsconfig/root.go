package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func Create(clientCertFile, clientKeyFile, caCertFile string) (tlsConfig *tls.Config, err error) {
	tlsConfig = &tls.Config{InsecureSkipVerify: false}
	if clientCertFile != "" {
		clientCert, err1 := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err1 != nil {
			return nil, fmt.Errorf("failed to load client cert/key: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{clientCert}
	}
	if caCertFile != "" {
		caCert, err := os.ReadFile(caCertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA cert: %w", err)
		}
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			return nil, fmt.Errorf("failed to append CA cert to pool")
		}
		tlsConfig.RootCAs = caCertPool
	}
	return tlsConfig, nil
}
