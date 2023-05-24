package fatturapa

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/invopop/gobl"
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
	Certificate   *xmldsig.Certificate
	WithTimestamp bool
	Transmitter   *Transmitter
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

// NewConverter returns a new GOBL to XML Converter with the given options
func NewConverter(opts ...Option) *Converter {
	c := new(Converter)
	c.Config = new(Config)
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// UnmarshalGOBL converts the given JSON document to a GOBL Envelope
func UnmarshalGOBL(reader io.Reader) (*gobl.Envelope, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, err
	}

	env := new(gobl.Envelope)
	if err := json.Unmarshal(buf.Bytes(), env); err != nil {
		return nil, err
	}

	return env, nil
}
