package fatturapa

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/xmldsig"
)

// <p:FatturaElettronica xmlns:ds="http://www.w3.org/2000/09/xmldsig#"
//   xmlns:p="http://ivaservizi.agenziaentrate.gov.it/docs/xsd/fatture/v1.2"
//   xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" versione="FPA12" xsi:schemaLocation="http://ivaservizi.agenziaentrate.gov.it/docs/xsd/fatture/v1.2 http://www.fatturapa.gov.it/export/fatturazione/sdi/fatturapa/v1.2/Schema_del_file_xml_FatturaPA_versione_1.2.xsd">

// Namespace used for FatturaPA. DSig stuff is handled in the signatures.
const (
	NamespaceFatturaPA = "http://ivaservizi.agenziaentrate.gov.it/docs/xsd/fatture/v1.2"
	NamespaceDSig      = "http://www.w3.org/2000/09/xmldsig#"
	NamespaceXSI       = "http://www.w3.org/2001/XMLSchema-instance"
	Versione           = "FPA12"
	SchemaLocation     = "http://ivaservizi.agenziaentrate.gov.it/docs/xsd/fatture/v1.2 https://www.fatturapa.gov.it/export/documenti/fatturapa/v1.2.2/Schema_del_file_xml_FatturaPA_v1.2.2.xsd"
)

// Document is a pseudo-model for containing the XML document being created.
type Document struct {
	env     *gobl.Envelope `xml:"-"` // Envelope to convert.
	invoice *bill.Invoice  `xml:"-"` // Invoice contained in envelope.

	XMLName        xml.Name `xml:"p:FatturaElettronica"`
	FPANamespace   string   `xml:"xmlns:p,attr"`
	DSigNamespace  string   `xml:"xmlns:ds,attr"`
	XSINamespace   string   `xml:"xmlns:xsi,attr"`
	Versione       string   `xml:"versione,attr"`
	SchemaLocation string   `xml:"xsi:schemaLocation,attr"`

	FatturaElettronicaHeader *FatturaElettronicaHeader
	FatturaElettronicaBody   []*FatturaElettronicaBody

	Signature *xmldsig.Signature `xml:"ds:Signature,omitempty"`
}

// LoadGOBL will build a FatturaPA Document from the source buffer
func (c *Client) LoadGOBL(src io.Reader) (*Document, error) {
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

// NewInvoice expects the base envelope and provides a new Document
// containing the XML version.
func (c *Client) NewInvoice(env *gobl.Envelope) (*Document, error) {
	invoice, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, errors.New("expected an invoice")
	}

	// Make sure we're dealing with raw data
	invoice = invoice.RemoveIncludedTaxes(2)

	header, err := newFatturaElettronicaHeader(invoice, c, env.Head.UUID.String())
	if err != nil {
		return nil, err
	}

	body, err := newFatturaElettronicaBody(invoice)
	if err != nil {
		return nil, err
	}

	// Basic document headers
	d := &Document{
		env:                      env,
		invoice:                  invoice,
		FPANamespace:             NamespaceFatturaPA,
		DSigNamespace:            NamespaceDSig,
		XSINamespace:             NamespaceXSI,
		Versione:                 header.DatiTrasmissione.FormatoTrasmissione,
		SchemaLocation:           SchemaLocation,
		FatturaElettronicaHeader: header,
		FatturaElettronicaBody:   []*FatturaElettronicaBody{body},
	}

	return d, nil
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
	// data, err := xml.MarshalIndent(d, "", "  ") // not compatible with certificates
	data, err := xml.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("marshal document: %w", err)
	}
	if _, err := buf.Write(data); err != nil {
		return nil, fmt.Errorf("writing to buffer: %w", err)
	}
	return buf, nil
}
