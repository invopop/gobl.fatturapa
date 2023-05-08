package fatturapa

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/invopop/gobl"
	"github.com/invopop/xmldsig"
)

// COnverter contains information related to the entity using this library
// to submit invoices to SDI.
type Converter struct {
	Transmitter *Transmitter
	Config      *Config
}

type Transmitter struct {
	CountryCode string
	TaxID       string
}

type Config struct {
	Certificate   *xmldsig.Certificate
	WithTimestamp bool
}

type Option func(*Converter)

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

func NewConverter(transmitter *Transmitter, opts ...Option) *Converter {
	c := new(Converter)
	c.Config = new(Config)
	c.Transmitter = transmitter
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// LoadGOBL will build a FatturaPA Document from the source buffer
func (c *Converter) LoadGOBL(src io.Reader) (*Document, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		return nil, err
	}

	env := new(gobl.Envelope)
	if err := json.Unmarshal(buf.Bytes(), env); err != nil {
		return nil, err
	}

	return c.NewInvoice(env)
}
