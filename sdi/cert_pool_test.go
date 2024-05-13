package sdi_test

import (
	"crypto/x509"
	"testing"

	sdi "github.com/invopop/gobl.fatturapa/sdi"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrepareCertPoolFromDir(t *testing.T) {
	t.Run("should return empty Cert Pool when folder has no cert files", func(t *testing.T) {
		expectedCertPool := x509.NewCertPool()
		caCertPool, err := sdi.PrepareCertPoolFromDir(".")

		require.NoError(t, err)
		assert.Equal(t, expectedCertPool, caCertPool)
		assert.Equal(t, 0, len(caCertPool.Subjects())) // nolint:staticcheck //
	})
}
