package sdi_test

import (
	"crypto/x509"
	"net/http"
	"testing"

	sdi "github.com/invopop/gobl.fatturapa/sdi"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient(t *testing.T) {
	t.Run("should return HTTP Client object", func(t *testing.T) {
		httpClient := sdi.HTTPClient(nil)

		assert.IsType(t, &http.Client{}, httpClient)
	})
	t.Run("should set cert pool in HTTP Client object", func(t *testing.T) {
		caCertPool := x509.NewCertPool()
		httpClient := sdi.HTTPClient(caCertPool)

		transport := httpClient.Transport.(*http.Transport)
		assert.IsType(t, caCertPool, transport.TLSClientConfig.RootCAs)
	})
}
