package sdi

import (
	"crypto/tls"
	"fmt"

	resty "github.com/go-resty/resty/v2"
)

// ClientOptFunc defines function for customizing the client
type ClientOptFunc func(*ClientOpts)

// ClientOpts defines the client parameters
type ClientOpts struct {
	Client        *resty.Client
	PrivateKeyPEM []byte
	PublicCertPEM []byte
	CACertPool    []byte
	Verbose       bool
}

func defaultClientOpts() ClientOpts {
	return ClientOpts{
		Client:        resty.New(),
		PrivateKeyPEM: nil,
		PublicCertPEM: nil,
		CACertPool:    nil,
		Verbose:       false,
	}
}

// Client defines http client
type Client struct {
	ClientOpts
}

// WithClient allows to customize the http client used for making the requests
func WithClient(client *resty.Client) ClientOptFunc {
	return func(o *ClientOpts) {
		o.Client = client
	}
}

// WithDebugMode uses a more verbose client
func WithDebugMode(verbose bool) ClientOptFunc {
	return func(o *ClientOpts) {
		o.Client.SetDebug(verbose)
		o.Verbose = verbose
	}
}

// WithCACertPool allows customizing CA Certificates
func WithCACertPool(caCertPool []byte) ClientOptFunc {
	return func(o *ClientOpts) {
		o.CACertPool = caCertPool
		o.Client.SetRootCertificateFromString(string(caCertPool))
	}
}

// WithClientCertPair allows customizing client certificates for mutual TLS authentication
func WithClientCertPair(publicCertPEM, privateKeyPEM []byte) ClientOptFunc {
	return func(o *ClientOpts) {
		o.PublicCertPEM = publicCertPEM
		o.PrivateKeyPEM = privateKeyPEM
		if publicCertPEM != nil && privateKeyPEM != nil {
			cert, err := tls.X509KeyPair(publicCertPEM, privateKeyPEM)
			if err != nil {
				panic(fmt.Errorf("client certificate error: %s", err))
			}
			o.Client.SetCertificates(cert)
		}
	}
}

// NewClient returns a customizing client
func NewClient(opts ...ClientOptFunc) *Client {
	o := defaultClientOpts()
	for _, fn := range opts {
		fn(&o)
	}
	return &Client{
		ClientOpts: o,
	}
}
