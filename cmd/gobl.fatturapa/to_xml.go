// Package main implements the CLI as well as mage commands (toplevel mage.go)
package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/invopop/gobl"
	fatturapa "github.com/invopop/gobl.fatturapa"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/xmldsig"
	"github.com/spf13/cobra"
)

type toXMLOpts struct {
	*convertOpts
	cert          string
	password      string
	transmitter   string
	withTimestamp bool
	pretty        bool
}

func toXML(c *convertOpts) *toXMLOpts {
	return &toXMLOpts{convertOpts: c}
}

func (t *toXMLOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "to-xml [infile] [outfile]",
		Short: "Convert a GOBL JSON into a FatturaPA XML",
		Long:  `Convert from GOBL JSON format to FatturaPA XML format.`,
		RunE:  t.runE,
	}
	f := cmd.Flags()
	f.StringVarP(&t.cert, "cert", "c", "", "Certificate for signing in pkcs12 format")
	f.StringVarP(&t.password, "password", "x", "", "Password of the certificate")
	f.StringVarP(&t.transmitter, "transmitter", "T", "", "Tax ID of the transmitter. Must be prefixed by the country code")
	f.BoolVarP(&t.withTimestamp, "with-timestamp", "t", false, "Add timestamp to the output file")
	f.BoolVarP(&t.pretty, "pretty", "p", true, "Output pretty-printed XML")
	return cmd
}

func (t *toXMLOpts) runE(cmd *cobra.Command, args []string) error {
	input, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer input.Close() // nolint:errcheck

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(input); err != nil {
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

	outFile := outputFilename(args)

	converter, err := t.loadConverterFromConfig()
	if err != nil {
		return err
	}

	doc, err := converter.ConvertFromGOBL(env)
	if err != nil {
		return err
	}

	// Marshal to JSON
	var data []byte
	if t.pretty {
		data, err = xml.MarshalIndent(doc, "", "\t")
	} else {
		data, err = xml.Marshal(doc)
	}
	if err != nil {
		return fmt.Errorf("marshaling to JSON: %w", err)
	}

	if err = os.WriteFile(outFile, data, 0644); err != nil {
		return fmt.Errorf("writing fatturapa xml: %w", err)
	}

	return nil
}

func (t *toXMLOpts) loadConverterFromConfig() (*fatturapa.Converter, error) {
	var opts []fatturapa.Option

	if t.transmitter != "" {
		countryCode := t.transmitter[:2]
		taxID := t.transmitter[2:]

		code := l10n.TaxCountryCode(countryCode)
		err := code.Validate()
		if err != nil {
			return nil, fmt.Errorf("tax ID must be prefixed by a valid country code")
		}

		transmitter := fatturapa.Transmitter{
			CountryCode: countryCode,
			TaxID:       taxID,
		}

		opts = append(opts, fatturapa.WithTransmitterData(&transmitter))
	}

	if t.cert != "" {
		cert, err := t.loadCertificate()
		if err != nil {
			return nil, err
		}

		opts = append(opts, fatturapa.WithCertificate(cert))
	}

	if t.withTimestamp {
		opts = append(opts, fatturapa.WithTimestamp())
	}

	return fatturapa.NewConverter(
		opts...,
	), nil
}

func (t *toXMLOpts) loadCertificate() (*xmldsig.Certificate, error) {
	cert, err := xmldsig.LoadCertificate(t.cert, t.password)
	if err != nil {
		return nil, fmt.Errorf("loading certificate %s: %w", t.cert, err)
	}

	return cert, nil
}
