package fatturapa

import (
	"github.com/invopop/xmldsig"
)

// Client contains information related to the entity using this library
// to submit invoices to SDI.
type Client struct {
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

type Option func(*Client)

// WithCertificate will ensure the XML document is signed with the given certificate
func WithCertificate(cert *xmldsig.Certificate) Option {
	return func(c *Client) {
		c.Config.Certificate = cert
	}
}

// WithTimestamp will ensure the XML document is timestamped
func WithTimestamp() Option {
	return func(c *Client) {
		c.Config.WithTimestamp = true
	}
}

func NewClient(transmitter *Transmitter, opts ...Option) *Client {
	c := new(Client)
	c.Config = new(Config)
	c.Transmitter = transmitter
	for _, opt := range opts {
		opt(c)
	}

	return c
}
