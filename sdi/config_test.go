package sdi_test

import (
	"testing"

	sdi "github.com/invopop/gobl.fatturapa/sdi"
	"github.com/stretchr/testify/assert"
)

func TestConfigSettings(t *testing.T) {
	t.Run("should return settings for development configuration", func(t *testing.T) {
		config := sdi.DevelopmentSdIConfig

		assert.Equal(t, "development", config.Environment)
		assert.Equal(t, "http://localhost:8080", config.Host)
		assert.Equal(t, "http://localhost:8080/ricevi_file", config.SOAPReceiveFileEndpoint())
		assert.Equal(t, "http://localhost:8080/RicezioneFatture", config.SOAPReceiveInvoicesEndpoint())
		assert.Equal(t, "http://localhost:8080/ricevi_notifica", config.SOAPReceiveNotificationEndpoint())
		assert.Equal(t, "http://localhost:8080/TrasmissioneFatture", config.SOAPTransmitInvoicesEndpoint())
	})

	t.Run("should return settings for production configuration", func(t *testing.T) {
		config := sdi.ProductionSdIConfig

		assert.Equal(t, "production", config.Environment)
		assert.Equal(t, "https://servizi.fatturapa.it", config.Host)
	})

	t.Run("should return settings for test configuration", func(t *testing.T) {
		config := sdi.TestSdIConfig

		assert.Equal(t, "test", config.Environment)
		assert.Equal(t, "https://testservizi.fatturapa.it", config.Host)
	})
}
