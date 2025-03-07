// Package test provides tools for testing the library both manually as well as
// helpers for writing test code.
package test

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/invopop/gobl"
	fatturapa "github.com/invopop/gobl.fatturapa"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/xmldsig"
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
)

const (
	certificateFile     = "test.p12"
	certificatePassword = "invopop"
	PathGOBLFatturaPA   = "gobl.fatturapa"
	PathFatturaPAGOBL   = "fatturapa.gobl"
)

// UpdateOut is a flag that can be set to update example files
var UpdateOut = flag.Bool("update", false, "Update the example files in test/data and test/data/out")

// NewConverter returns a fatturapa.Converter with the test certificate and
// transmitter data.
func NewConverter() *fatturapa.Converter {
	cert, err := loadCertificate()

	if err != nil {
		panic(err)
	}

	transmitter := &fatturapa.Transmitter{
		CountryCode: string(l10n.IT),
		TaxID:       "01234567890",
	}

	// Set a fixed time to get deterministic signatures
	ts, _ := time.Parse(time.RFC3339, "2022-02-01T04:00:00Z")

	converter := fatturapa.NewConverter(
		fatturapa.WithTransmitterData(transmitter),
		fatturapa.WithCertificate(cert),
		fatturapa.WithCurrentTime(ts),
	)

	return converter
}

// ConvertFromGOBL takes the GOBL test data and converts into XML
func ConvertFromGOBL(env *gobl.Envelope, converter ...*fatturapa.Converter) (*fatturapa.Document, error) {
	var c *fatturapa.Converter

	if len(converter) == 0 {
		c = NewConverter()
	} else {
		c = converter[0]
	}

	doc, err := c.ConvertFromGOBL(env)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// ConvertToGOBL takes the XML test data and converts into a GOBL envelope
func ConvertToGOBL(doc []byte, converter ...*fatturapa.Converter) (*gobl.Envelope, error) {
	var c *fatturapa.Converter

	if len(converter) == 0 {
		c = NewConverter()
	} else {
		c = converter[0]
	}

	env, err := c.ConvertToGOBL(doc)
	if err != nil {
		return nil, err
	}

	return env, nil
}

// GetDataPath returns the path where test can find data files
// to be used in tests
func GetDataPath() string {
	return getRootFolder() + "/test/data/"
}

// ModifyInvoice takes a GOBL envelope and modifies the invoice
func ModifyInvoice(env *gobl.Envelope, modifyFunc func(*bill.Invoice)) {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		panic("error extracting invoice")
	}

	modifyFunc(inv)

	doc, err := schema.NewObject(inv)
	if err != nil {
		panic(err)
	}

	env.Document = doc
}

// LoadTestFile loads a test file from the test/data folder as a GOBL envelope
func LoadTestFile(file string) *gobl.Envelope {
	path := filepath.Join(GetDataPath(), file)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(f); err != nil {
		panic(err)
	}

	out, err := gobl.Parse(buf.Bytes())
	if err != nil {
		panic(err)
	}

	var env *gobl.Envelope
	switch doc := out.(type) {
	case *gobl.Envelope:
		env = doc
	default:
		env = gobl.NewEnvelope()
		if err := env.Insert(doc); err != nil {
			panic(err)
		}
	}

	if err := env.Calculate(); err != nil {
		panic(err)
	}

	if err := env.Validate(); err != nil {
		panic(err)
	}

	if *UpdateOut {
		data, err := json.MarshalIndent(env, "", "\t")
		if err != nil {
			panic(err)
		}

		if err := os.WriteFile(path, data, 0644); err != nil {
			panic(err)
		}
	}

	return env
}

// LoadSchema loads a XSD schema for validating XML documents
func LoadSchema() (*xsd.Schema, error) {
	schemaPath := filepath.Join("schemas", "FatturaPA_v1.2.2.xsd")
	schema, err := xsd.ParseFromFile(schemaPath)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

// ValidateXML validates an XML document against a XSD schema
func ValidateXML(schema *xsd.Schema, doc []byte) []error {
	xmlDoc, err := libxml2.ParseString(string(doc))
	if err != nil {
		return []error{err}
	}

	err = schema.Validate(xmlDoc)
	if err != nil {
		return err.(xsd.SchemaValidationError).Errors()
	}

	return nil
}

func loadCertificate() (*xmldsig.Certificate, error) {
	certificatesPath := getRootFolder() + "/test/certificates/"

	f := path.Join(certificatesPath, certificateFile)
	return xmldsig.LoadCertificate(f, certificatePassword)
}

func getRootFolder() string {
	cwd, _ := os.Getwd()

	for !isRootFolder(cwd) {
		cwd = removeLastEntry(cwd)
	}

	return cwd
}

func isRootFolder(dir string) bool {
	files, _ := os.ReadDir(dir)

	for _, file := range files {
		if file.Name() == "go.mod" {
			return true
		}
	}

	return false
}

func removeLastEntry(dir string) string {
	lastEntry := "/" + filepath.Base(dir)
	i := strings.LastIndex(dir, lastEntry)
	return dir[:i]
}
