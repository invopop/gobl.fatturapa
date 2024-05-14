package sdi

// Config represents the configuration structure
type Config struct {
	Environment string
	Host        string
}

// SOAPReceiveFileEndpoint returns an endpoint for receive file via SOAP
func (c Config) SOAPReceiveFileEndpoint() string {
	return c.Host + "/ricevi_file"
}

// SOAPReceiveInvoicesEndpoint returns an endpoint for receive invoices via SOAP
func (c Config) SOAPReceiveInvoicesEndpoint() string {
	return c.Host + "/RicezioneFatture"
}

// SOAPReceiveNotificationEndpoint returns an endpoint for receive invoices via SOAP
func (c Config) SOAPReceiveNotificationEndpoint() string {
	return c.Host + "/ricevi_notifica"
}

// SOAPTransmitInvoicesEndpoint returns an endpoint for transmit invoices via SOAP
func (c Config) SOAPTransmitInvoicesEndpoint() string {
	return c.Host + "/TrasmissioneFatture"
}

var (
	// DevelopmentSdIConfig is the settings for development
	DevelopmentSdIConfig = Config{
		Environment: "development",
		Host:        "http://localhost:8080",
	}

	// ProductionSdIConfig is the settings for production
	ProductionSdIConfig = Config{
		Environment: "production",
		Host:        "https://servizi.fatturapa.it",
	}

	// TestSdIConfig is the settings for test
	TestSdIConfig = Config{
		Environment: "test",
		Host:        "https://testservizi.fatturapa.it",
	}
)
