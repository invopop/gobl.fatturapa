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

func TestConverter() *fatturapa.Converter {
	cert, err := loadCertificate()

	if err != nil {
		panic(err)
	}

	transmitter := &fatturapa.Transmitter{
		CountryCode: "IT",
		TaxID:       "01234567890",
	}

	converter := fatturapa.NewConverter(
		fatturapa.WithTransmissionData(transmitter),
		fatturapa.WithCertificate(cert),
	)

	return converter
}

// LoadCertificate will return the standard test certificate
func loadCertificate() (*xmldsig.Certificate, error) {
	certificatesPath := getRootFolder() + "/test/certificates/"

	f := path.Join(certificatesPath, certificateFile)
	return xmldsig.LoadCertificate(f, certificatePassword)
}
