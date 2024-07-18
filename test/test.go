// Package test provides tools for testing the library both manually as well as
// helpers for writing test code.
package test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/invopop/gobl"
	fatturapa "github.com/invopop/gobl.fatturapa"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/xmldsig"
)

const (
	certificateFile     = "test.p12"
	certificatePassword = "invopop"
)

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

	converter := fatturapa.NewConverter(
		fatturapa.WithTransmitterData(transmitter),
		fatturapa.WithCertificate(cert),
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

// TestConversion takes the .json invoices generated previously and converts them
// into XML fatturapa documents.
func TestConversion() error { // nolint:revive
	var files []string
	err := filepath.Walk(GetDataPath(), func(path string, _ os.FileInfo, _ error) error {
		if filepath.Ext(path) == ".json" {
			files = append(files, filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Printf("processing file: %v\n", file)

		envelopeReader, err := os.Open(GetDataPath() + file)
		if err != nil {
			return err
		}

		env, err := fatturapa.UnmarshalGOBL(envelopeReader)
		if err != nil {
			return err
		}

		doc, err := ConvertFromGOBL(env, NewConverter())
		if err != nil {
			return err
		}

		data, err := doc.Bytes()
		if err != nil {
			return fmt.Errorf("extracting document bytes: %w", err)
		}

		np := strings.TrimSuffix(file, filepath.Ext(file)) + ".xml"
		err = os.WriteFile(GetDataPath()+"/"+np, data, 0644)
		if err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
	}

	return nil
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
	f, err := os.Open(GetDataPath() + file)
	if err != nil {
		panic(err)
	}

	env, err := fatturapa.UnmarshalGOBL(f)
	if err != nil {
		panic(err)
	}

	return env
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
