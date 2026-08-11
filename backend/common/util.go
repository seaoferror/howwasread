package common

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func FetchPublicWANIP() (net.IP, error) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public IP: %w", err)
	}
	defer resp.Body.Close()

	ipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ipStr := strings.TrimSpace(string(ipBytes))
	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil {
		return nil, fmt.Errorf("invalid IP format received: %s", ipStr)
	}

	return parsedIP, nil
}

func CreateTlSConfig(clientCertFile, clientKeyFile, caCertFile string) (tlsConfig *tls.Config, err error) {
	tlsConfig = &tls.Config{InsecureSkipVerify: false}
	if clientCertFile != "" {
		clientCert, err1 := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err1 != nil {
			return nil, fmt.Errorf("failed to load client cert/key: %w", err1)
		}
		tlsConfig.Certificates = []tls.Certificate{clientCert}
	}
	if caCertFile != "" {
		caCert, err2 := os.ReadFile(caCertFile)
		if err2 != nil {
			return nil, fmt.Errorf("failed to read CA cert: %w", err2)
		}
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			return nil, fmt.Errorf("failed to append CA cert to pool")
		}
		tlsConfig.RootCAs = caCertPool
	}
	return tlsConfig, nil
}
