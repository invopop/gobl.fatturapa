package test

import (
	"path"

	fatturapa "github.com/invopop/gobl.fatturapa"
	"github.com/invopop/xmldsig"
)

const (
	certificateFile     = "test2.p12"
	certificatePassword = "invopop"
)

func TestClient() *fatturapa.Client {
	cert, err := loadCertificate()

	if err != nil {
		panic(err)
	}

	client := fatturapa.NewClient(
		&fatturapa.Transmitter{
			CountryCode: "IT",
			TaxID:       "01234567890",
		},
		fatturapa.WithCertificate(cert),
	)

	return client
}

// LoadCertificate will return the standard test certificate
func loadCertificate() (*xmldsig.Certificate, error) {
	certificatesPath := getRootFolder() + "/test/certificates/"

	f := path.Join(certificatesPath, certificateFile)
	return xmldsig.LoadCertificate(f, certificatePassword)
}
