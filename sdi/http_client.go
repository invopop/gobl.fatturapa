package sdi

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"
)

// HTTPClient returns HTTP client with Cert Pool
func HTTPClient(caCertPool *x509.CertPool) *http.Client {
	httpClient := &http.Client{
		Timeout: 2500 * time.Millisecond,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	return httpClient
}
