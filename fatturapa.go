// Package fatturapa implements the conversion from GOBL to FatturaPA XML.
package fatturapa

import (
	"bytes"
	"errors"
	"fmt"

	// Momentary fix for nbio/xml library not properly handling namespace attributes but still needed for unmarshalling
	"encoding/xml"

	nbioXML "github.com/nbio/xml"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/it/sdi"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/xmldsig"
)

// <p:FatturaElettronica xmlns:ds="http://www.w3.org/2000/09/xmldsig#"
//   xmlns:p="http://ivaservizi.agenziaentrate.gov.it/docs/xsd/fatture/v1.2"
//   xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" versione="FPA12" xsi:schemaLocation="http://ivaservizi.agenziaentrate.gov.it/docs/xsd/fatture/v1.2 http://www.fatturapa.gov.it/export/fatturazione/sdi/fatturapa/v1.2/Schema_del_file_xml_FatturaPA_versione_1.2.xsd">

// Namespace used for FatturaPA. DSig stuff is handled in the signatures.
const (
	namespaceFatturaPA = "http://ivaservizi.agenziaentrate.gov.it/docs/xsd/fatture/v1.2"
	namespaceDSig      = "http://www.w3.org/2000/09/xmldsig#"
	namespaceXSI       = "http://www.w3.org/2001/XMLSchema-instance"
	schemaLocation     = "http://ivaservizi.agenziaentrate.gov.it/docs/xsd/fatture/v1.2 https://www.fatturapa.gov.it/export/documenti/fatturapa/v1.2.2/Schema_del_file_xml_FatturaPA_v1.2.2.xsd"
)

// Document is a pseudo-model for containing the XML document being created.
type Document struct {
	env *gobl.Envelope `xml:"-"` // Envelope to convert.

	XMLName        xml.Name `xml:"p:FatturaElettronica"`
	FPANamespace   string   `xml:"xmlns:p,attr"`
	DSigNamespace  string   `xml:"xmlns:ds,attr"`
	XSINamespace   string   `xml:"xmlns:xsi,attr"`
	Versione       string   `xml:"versione,attr"`
	SchemaLocation string   `xml:"xsi:schemaLocation,attr"`

	Header *Header `xml:"FatturaElettronicaHeader"`
	Body   []*Body `xml:"FatturaElettronicaBody"`

	Signature *xmldsig.Signature `xml:"ds:Signature,omitempty"`
}

// ConvertFromGOBL expects the base envelope and provides a new Document
// containing the XML version.
func (c *Converter) ConvertFromGOBL(env *gobl.Envelope) (*Document, error) {
	invoice, ok := env.Extract().(*bill.Invoice)
	if !ok || invoice == nil {
		return nil, errors.New("expected an invoice")
	}

	// Make sure we're dealing with raw data
	if err := invoice.RemoveIncludedTaxes(); err != nil {
		return nil, err
	}

	TransmissionData := c.newTransmissionData(invoice, env)

	header := newHeader(invoice, TransmissionData)

	body, err := newBody(invoice)
	if err != nil {
		return nil, err
	}

	// Basic document headers
	d := &Document{
		env:            env,
		FPANamespace:   namespaceFatturaPA,
		DSigNamespace:  namespaceDSig,
		XSINamespace:   namespaceXSI,
		Versione:       formatoTransmissione(invoice),
		SchemaLocation: schemaLocation,
		Header:         header,
		Body:           []*Body{body},
	}

	if c.Config.Certificate != nil {
		if err := d.sign(c.Config); err != nil {
			return nil, err
		}
	}

	return d, nil
}

// ConvertToGOBL expects the XML document bytes and provides a new GOBL
// envelope containing the invoice.
func (c *Converter) ConvertToGOBL(doc []byte) (*gobl.Envelope, error) {
	d := &Document{}
	if err := nbioXML.Unmarshal(doc, d); err != nil {
		return nil, fmt.Errorf("unmarshal document: %w", err)
	}

	// Verify signature. Standin for now.
	// Skip signature verification for now
	/*
		if d.Signature == nil {
			return nil, errors.New("signature is missing")
		}
	*/

	// Create a new invoice with empty fields so that converter can fill it
	inv := new(bill.Invoice)

	inv.Addons = tax.WithAddons(sdi.V1)

	// Retrieves information from the header and adds it to the invoice
	goblBillInvoiceAddHeader(inv, d.Header)

	// Retrieves information from the body and adds it to the invoice
	// TODO: add support for multiple bodies
	err := goblBillInvoiceAddBody(inv, d.Body[0])
	if err != nil {
		return nil, err
	}

	// Generate envelope
	env, err := gobl.Envelop(inv)
	if err != nil {
		return nil, err
	}

	if err := env.Validate(); err != nil {
		return nil, err
	}

	return env, nil
}

// Buffer returns a byte buffer representation of the complete XML document.
func (d *Document) Buffer() (*bytes.Buffer, error) {
	return d.buffer(xml.Header)
}

// String converts a struct representation to its string representation
func (d *Document) String() (string, error) {
	buf, err := d.Buffer()
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Bytes returns the XML document bytes
func (d *Document) Bytes() ([]byte, error) {
	buf, err := d.Buffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (d *Document) buffer(base string) (*bytes.Buffer, error) {
	buf := bytes.NewBufferString(base)
	//data, err := xml.MarshalIndent(d, "", "  ")
	data, err := xml.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("marshal document: %w", err)
	}
	if _, err := buf.Write(data); err != nil {
		return nil, fmt.Errorf("writing to buffer: %w", err)
	}
	return buf, nil
}
