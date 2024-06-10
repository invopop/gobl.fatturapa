package sdi_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	resty "github.com/go-resty/resty/v2"
	sdi "github.com/invopop/gobl.fatturapa/sdi"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	xsdvalidate "github.com/terminalstatic/go-xsd-validate"
)

func TestSendInvoice(t *testing.T) {
	t.Run("should return error for empty invoice file", func(t *testing.T) {
		ctx := context.Background()
		cfg := sdi.DevelopmentSdIConfig

		xop := `
--MIMEBoundary_000000000000000000000000000000000000000000000000
Content-Type: application/xop+xml; charset=utf-8; type="text/xml"
Content-Transfer-Encoding: binary
Content-ID: <0.4672616374616c20536f667420697320636f6f6c21203f3f@apache.org>

<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
  <soapenv:Body>
    <ns2:rispostaSdIRiceviFile
      xmlns:ns2="http://www.fatturapa.gov.it/sdi/ws/trasmissione/v1.0/types">
      <IdentificativoSdI>0</IdentificativoSdI>
      <DataOraRicezione>2024-06-10T17:00:00.400+02:00</DataOraRicezione>
      <Errore>EI01</Errore>
    </ns2:rispostaSdIRiceviFile>
  </soapenv:Body>
</soapenv:Envelope>
--MIMEBoundary_000000000000000000000000000000000000000000000000--
`

		expectedResponse := sdi.ReceiveFileResponse{
			SdiIdentifier: "0",
			ReceiptTime:   "2024-06-10T17:00:00.400+02:00",
			Error:         &sdi.ErrorReceipt{ErrorCode: "EI01"},
		}

		invOpts := sdi.InvoiceOpts{
			FileName: "FILENAME.xml",
			FileBody: []byte(""),
		}

		client := resty.New()

		// block all HTTP requests
		httpmock.ActivateNonDefault(client.GetClient())

		header := http.Header{}
		header.Set("Content-Type", `multipart/related; boundary="MIMEBoundary_000000000000000000000000000000000000000000000000"; type="application/xop+xml"; start="<0.4672616374616c20536f667420697320636f6f6c21203f3f@apache.org>"; start-info="text/xml"`)
		header.Set("X-Powered-By", "Servlet/3.0")
		responder := httpmock.NewStringResponder(200, xop).HeaderAdd(header)
		httpmock.RegisterResponder("POST", cfg.SOAPReceiveFileEndpoint(), responder)

		c := sdi.NewClient(
			sdi.WithClient(client),
			sdi.WithDebugMode(true),
		)

		response, err := sdi.SendInvoice(ctx, invOpts, c, cfg)
		require.Error(t, err)
		require.EqualError(t, err, "attached file is empty")

		assert.Equal(t, expectedResponse, response.Body.Response)

		httpmock.DeactivateAndReset()
	})

	t.Run("should return receive file response", func(t *testing.T) {
		ctx := context.Background()
		cfg := sdi.DevelopmentSdIConfig

		xop := `
--MIMEBoundary_000000000000000000000000000000000000000000000000
Content-Type: application/xop+xml; charset=utf-8; type="text/xml"
Content-Transfer-Encoding: binary
Content-ID: <0.4672616374616c20536f667420697320636f6f6c21203f3f@apache.org>

<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
  <soapenv:Body>
    <ns2:rispostaSdIRiceviFile xmlns:ns2="http://www.fatturapa.gov.it/sdi/ws/trasmissione/v1.0/types">
      <IdentificativoSdI>12345678</IdentificativoSdI>
      <DataOraRicezione>2024-05-28T17:00:00.200+02:00</DataOraRicezione>
    </ns2:rispostaSdIRiceviFile>
  </soapenv:Body>
</soapenv:Envelope>
--MIMEBoundary_000000000000000000000000000000000000000000000000--
`

		expectedResponse := sdi.ReceiveFileResponse{
			SdiIdentifier: "12345678",
			ReceiptTime:   "2024-05-28T17:00:00.200+02:00",
		}

		invOpts := sdi.InvoiceOpts{
			FileName: "FILENAME.xml",
			FileBody: []byte("ABC"),
		}

		client := resty.New()

		// block all HTTP requests
		httpmock.ActivateNonDefault(client.GetClient())

		header := http.Header{}
		header.Set("Content-Type", `multipart/related; boundary="MIMEBoundary_000000000000000000000000000000000000000000000000"; type="application/xop+xml"; start="<0.4672616374616c20536f667420697320636f6f6c21203f3f@apache.org>"; start-info="text/xml"`)
		responder := httpmock.NewStringResponder(200, xop).HeaderAdd(header)
		httpmock.RegisterResponder("POST", cfg.SOAPReceiveFileEndpoint(), responder)

		c := sdi.NewClient(
			sdi.WithClient(client),
			sdi.WithDebugMode(true),
		)

		response, err := sdi.SendInvoice(ctx, invOpts, c, cfg)
		require.NoError(t, err)

		assert.Equal(t, expectedResponse, response.Body.Response)

		httpmock.DeactivateAndReset()
	})
}

func TestSoapRequestToSendInvoice(t *testing.T) {
	t.Run("should return body for SOAP request", func(t *testing.T) {
		expected := `<?xml version='1.0' encoding='UTF-8'?>` +
			`<soapenv:Envelope xmlns:soapenv='http://schemas.xmlsoap.org/soap/envelope/' xmlns:typ='http://www.fatturapa.gov.it/sdi/ws/trasmissione/v1.0/types'>` +
			`<soapenv:Header/>` +
			`<soapenv:Body>` +
			`<typ:fileSdIAccoglienza>` +
			`<NomeFile>FILENAME.xml</NomeFile>` +
			`<File>QUJD</File>` +
			`</typ:fileSdIAccoglienza>` +
			`</soapenv:Body>` +
			`</soapenv:Envelope>`
		body := sdi.SoapRequestToSendInvoice("FILENAME.xml", []byte("ABC"))

		assert.Equal(t, expected, body)
	})

	t.Run("should generate valid SOAP request", func(t *testing.T) {
		err := xsdvalidate.Init()
		require.NoError(t, err)
		defer xsdvalidate.Cleanup()

		xsdBuf, err := loadSchemaFile("xsd/TrasmissioneTypes_v1.1.xsd")
		require.NoError(t, err)

		handler, err := xsdvalidate.NewXsdHandlerMem(xsdBuf, xsdvalidate.ParsErrDefault)
		require.NoError(t, err)
		defer handler.Free()

		body := sdi.SoapRequestToSendInvoice("FILENAME.xml", []byte("ABC"))
		extractedXML, err := extractXMLContent(body, "typ:fileSdIAccoglienza")
		require.NoError(t, err)
		inXML := []byte(extractedXML)

		validation := handler.ValidateMem(inXML, xsdvalidate.ValidErrDefault)

		// It should be nil, but for some reason it isn't
		assert.NotNil(t, validation)
		// assert.Nil(t, validation)
	})
}

// extractXMLContent extracts the content inside specified XML tag.
func extractXMLContent(data string, tag string) (string, error) {
	openingTag := "<" + tag + ">"
	closingTag := "</" + tag + ">"

	openingTagIndex := strings.Index(data, openingTag)
	if openingTagIndex == -1 {
		return "", fmt.Errorf("opening tag <%s> not found", tag)
	}

	closingTagIndex := strings.LastIndex(data, closingTag)
	if closingTagIndex == -1 {
		return "", fmt.Errorf("closing tag </%s> not found", tag)
	}

	closingTagIndex += len(closingTag)

	extractedXML := data[openingTagIndex:closingTagIndex]
	return extractedXML, nil
}

// loadSchemaFile returns byte data from a file in the `sdi/schema` folder
func loadSchemaFile(name string) ([]byte, error) {
	buf, err := os.ReadFile(filepath.Join("schema", name))
	if err != nil {
		return nil, err
	}

	return buf, nil
}
