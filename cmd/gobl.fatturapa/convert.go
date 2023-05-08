package main

import (
	"fmt"

	fatturapa "github.com/invopop/gobl.fatturapa"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/xmldsig"
	"github.com/spf13/cobra"
)

type convertOpts struct {
	*rootOpts
	cert          string
	password      string
	taxID         string
	withTimestamp bool
}

func convert(o *rootOpts) *convertOpts {
	return &convertOpts{rootOpts: o}
}

func (c *convertOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert [infile] [outfile]",
		Short: "Convert a GOBL JSON into a FatturaPA XML",
		RunE:  c.runE,
	}
	f := cmd.Flags()
	f.StringVarP(&c.cert, "cert", "c", "", "Certificate for signing in pkcs12 format")
	f.StringVarP(&c.password, "password", "p", "", "Password of the certificate")
	f.StringVarP(&c.taxID, "tax-id", "x", "", "Tax ID of the transmitter. Must be prefixed by the country code")
	f.BoolVarP(&c.withTimestamp, "with-timestamp", "t", false, "Add timestamp to the output file")

	return cmd
}

func (c *convertOpts) runE(cmd *cobra.Command, args []string) error {
	// ctx := commandContext(cmd)

	input, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer input.Close() // nolint:errcheck

	out, err := c.openOutput(cmd, args)
	if err != nil {
		return err
	}
	defer out.Close() // nolint:errcheck

	converter, err := loadConverterFromConfig(c)
	if err != nil {
		return err
	}

	doc, err := converter.LoadGOBL(input)
	if err != nil {
		return err
	}

	data, err := doc.Bytes()
	if err != nil {
		return fmt.Errorf("generating fatturapa xml: %w", err)
	}

	if _, err = out.Write(data); err != nil {
		return fmt.Errorf("writing fatturapa xml: %w", err)
	}

	return nil
}

func loadConverterFromConfig(c *convertOpts) (*fatturapa.Converter, error) {
	var opts []fatturapa.Option

	if c.taxID != "" {
		countryCode := c.taxID[:2]
		taxID := c.taxID[2:]

		code := l10n.CountryCode(countryCode)
		err := code.Validate()
		if err != nil {
			return nil, fmt.Errorf("tax ID must be prefixed by a valid country code")
		}

		transmitter := fatturapa.Transmitter{
			CountryCode: countryCode,
			TaxID:       taxID,
		}

		opts = append(opts, fatturapa.WithTransmissionData(&transmitter))
	}

	if c.cert != "" {
		cert, err := loadCertificate(c.cert, c.password)
		if err != nil {
			return nil, err
		}

		opts = append(opts, fatturapa.WithCertificate(cert))
	}

	if c.withTimestamp {
		opts = append(opts, fatturapa.WithTimestamp())
	}

	return fatturapa.NewConverter(
		opts...,
	), nil
}

func loadCertificate(certPath, password string) (*xmldsig.Certificate, error) {
	cert, err := xmldsig.LoadCertificate(certPath, password)
	if err != nil {
		return nil, fmt.Errorf("loading certificate %s: %w", certPath, err)
	}

	return cert, nil
}
