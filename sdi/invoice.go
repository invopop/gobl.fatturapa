package sdi

import (
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"
)

// SendInvoiceResponse defines the post invoice response structure
type SendInvoiceResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		Response ReceiveFileResponse `xml:"rispostaSdIRiceviFile"`
	} `xml:"Body"`
}

// InvoiceOpts defines the send invoice parameters
type InvoiceOpts struct {
	FileName string
	FileBody []byte
}

// SendInvoice sends invoice content to SdI
func SendInvoice(ctx context.Context, invOpts InvoiceOpts, c *Client, cfg Config) (*SendInvoiceResponse, error) {
	soapEndpoint := cfg.SOAPReceiveFileEndpoint()

	body := SoapRequestToSendInvoice(invOpts.FileName, invOpts.FileBody)
	resp, err := c.Client.R().
		SetBody(body).
		SetContext(ctx).
		Post(soapEndpoint)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("response error: %v", resp)
	}

	response := &SendInvoiceResponse{}
	err = parseMultipartResponse(resp, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SoapRequestToSendInvoice prepares the request content for SOAP to send an invoice
func SoapRequestToSendInvoice(fileName string, fileBody []byte) string {
	blob := base64.StdEncoding.EncodeToString(fileBody)

	return `<?xml version='1.0' encoding='UTF-8'?>` +
		`<soapenv:Envelope xmlns:soapenv='http://schemas.xmlsoap.org/soap/envelope/' xmlns:typ='http://www.fatturapa.gov.it/sdi/ws/trasmissione/v1.0/types'>` +
		`<soapenv:Header/>` +
		`<soapenv:Body>` +
		`<typ:fileSdIAccoglienza>` +
		`<NomeFile>` + fileName + `</NomeFile>` +
		`<File>` + blob + `</File>` +
		`</typ:fileSdIAccoglienza>` +
		`</soapenv:Body>` +
		`</soapenv:Envelope>`
}
