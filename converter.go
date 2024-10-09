package fatturapa

import (
	"time"

	"github.com/invopop/xmldsig"
)

// Converter contains information related to the entity using this library
// to submit invoices to SDI.
type Converter struct {
	Config *Config
}

// Transmitter contains information about the entity integrating directly
// with the SDI to submit and receive invoices
type Transmitter struct {
	CountryCode string
	TaxID       string
}

// Config contains the configuration for the Converter
type Config struct {
	Certificate     *xmldsig.Certificate
	WithTimestamp   bool
	Transmitter     *Transmitter
	WithCurrentTime time.Time
}

// Option is a function that can be passed to NewConverter to configure it
type Option func(*Converter)

// WithTransmitterData will ensure the XML document contains the given transmitter data
func WithTransmitterData(transmitter *Transmitter) Option {
	return func(c *Converter) {
		c.Config.Transmitter = transmitter
	}
}

// WithCertificate will ensure the XML document is signed with the given certificate
func WithCertificate(cert *xmldsig.Certificate) Option {
	return func(c *Converter) {
		c.Config.Certificate = cert
	}
}

// WithTimestamp will ensure the XML document is timestamped
func WithTimestamp() Option {
	return func(c *Converter) {
		c.Config.WithTimestamp = true
	}
}

// WithCurrentTime will ensure the XML document is signed with the given current time
func WithCurrentTime(t time.Time) Option {
	return func(c *Converter) {
		c.Config.WithCurrentTime = t
	}
}

// NewConverter returns a new GOBL to XML Converter with the given options
func NewConverter(opts ...Option) *Converter {
	c := new(Converter)
	c.Config = new(Config)
	for _, opt := range opts {
		opt(c)
	}

	return c
}
