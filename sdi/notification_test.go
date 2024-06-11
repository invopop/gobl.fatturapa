package sdi_test

import (
	"fmt"
	"testing"

	sdi "github.com/invopop/gobl.fatturapa/sdi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseReceiptRejection(t *testing.T) {
	t.Run("should return error code 00102 and info about an invalid signature in parsed struct", func(t *testing.T) {
		//nolint:misspell
		xml := `
<?xml version="1.0" encoding="UTF-8"?>
<ns3:RicevutaScarto
  xmlns:ns3="http://ivaservizi.agenziaentrate.gov.it/docs/xsd/fattura/messaggi/v1.0"
  xmlns:ns2="http://www.w3.org/2000/09/xmldsig#" versione="1.0">
  <IdentificativoSdI>12345678</IdentificativoSdI>
  <NomeFile>IT01234567890_FPA01.xml</NomeFile>
  <Hash>4672616374616c20536f667420697320636f6f6c21203f3f3f3f3f3f3f3f3f3f</Hash>
  <DataOraRicezione>2024-06-11T16:00:00.000+00:00</DataOraRicezione>
  <ListaErrori>
    <Errore>
      <Codice>00102</Codice>
      <Descrizione>File non integro (firma non valida) : 00102&#13;</Descrizione>
      <Suggerimento>Verificare che il file sia firmato correttamente o che non sia stato modificato dopo l'apposizione della firma</Suggerimento>
    </Errore>
  </ListaErrori>
  <MessageId>100000000</MessageId>
  <ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="Signature1">
    <ds:SignedInfo>
      <ds:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#" />
      <ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256" />
      <ds:Reference Id="reference-document" URI="">
        <ds:Transforms>
          <ds:Transform Algorithm="http://www.w3.org/2002/06/xmldsig-filter2">
            <XPath xmlns="http://www.w3.org/2002/06/xmldsig-filter2" Filter="subtract">/descendant::ds:Signature</XPath>
          </ds:Transform>
        </ds:Transforms>
        <ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256" />
        <ds:DigestValue>jCGMAK8UHIv+RiBwCKdGsLRSpKkDExW24eN5tSF0gNg=</ds:DigestValue>
      </ds:Reference>
      <ds:Reference Id="reference-signedpropeties" Type="http://uri.etsi.org/01903#SignedProperties" URI="#SignedProperties_1">
        <ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256" />
        <ds:DigestValue>Z...0=</ds:DigestValue>
      </ds:Reference>
      <ds:Reference Id="reference-keyinfo" URI="#KeyInfoId">
        <ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256" />
        <ds:DigestValue>/+Q...E=</ds:DigestValue>
      </ds:Reference>
    </ds:SignedInfo>
    <ds:SignatureValue Id="SignatureValue1">SV...==</ds:SignatureValue>
    <ds:KeyInfo Id="KeyInfoId">
      <ds:X509Data>
        <ds:X509Certificate>MIIF...M5</ds:X509Certificate>
      </ds:X509Data>
    </ds:KeyInfo>
    <ds:Object>
      <xades:QualifyingProperties xmlns:xades="http://uri.etsi.org/01903/v1.3.2#" Target="#Signature1">
        <xades:SignedProperties Id="SignedProperties_1">
          <xades:SignedSignatureProperties>
            <xades:SigningTime>2024-06-11T16:00:00Z</xades:SigningTime>
          </xades:SignedSignatureProperties>
        </xades:SignedProperties>
      </xades:QualifyingProperties>
    </ds:Object>
  </ds:Signature>
</ns3:RicevutaScarto>
`

		output, err := sdi.ParseReceiptRejection([]byte(xml))
		require.Error(t, err)
		assert.Equal(t, fmt.Errorf("sdi error code list: [00102]"), err)

		errors := output.ErrorList.Error
		assert.Equal(t, 1, len(errors))
		assert.Equal(t, "00102", errors[0].Code)
	})
}
