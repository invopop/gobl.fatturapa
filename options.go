package fatturapa

import (
	"time"

	"github.com/invopop/xmldsig"
)

// Transmitter contains information about the entity integrating directly
// with the SDI to submit and receive invoices
type Transmitter struct {
	CountryCode string
	TaxID       string
}

// Config contains the configuration for the Converter
type config struct {
	Certificate     *xmldsig.Certificate
	WithTimestamp   bool
	Transmitter     *Transmitter
	WithCurrentTime time.Time
}

// Option is a function that can be passed to NewConverter to configure it
type Option func(*config)

// WithTransmitterData will ensure the XML document contains the given transmitter data
func WithTransmitterData(transmitter *Transmitter) Option {
	return func(c *config) {
		c.Transmitter = transmitter
	}
}

// WithCertificate will ensure the XML document is signed with the given certificate
func WithCertificate(cert *xmldsig.Certificate) Option {
	return func(c *config) {
		c.Certificate = cert
	}
}

// WithTimestamp will ensure the XML document is timestamped
func WithTimestamp() Option {
	return func(c *config) {
		c.WithTimestamp = true
	}
}

// WithCurrentTime will ensure the XML document is signed with the given current time
func WithCurrentTime(t time.Time) Option {
	return func(c *config) {
		c.WithCurrentTime = t
	}
}

// NewConverter returns a new GOBL to XML Converter with the given options
func parseOptions(opts ...Option) *config {
	c := new(config)
	for _, opt := range opts {
		opt(c)
	}

	return c
}
