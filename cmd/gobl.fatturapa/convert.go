package main

import (
	"fmt"
	"log"
	"os"

	fatturapa "github.com/invopop/gobl.fatturapa"
	"github.com/invopop/xmldsig"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

type convertOpts struct {
	*rootOpts
	cert          string
	password      string
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

	client, err := loadClientFromConfig(c)
	if err != nil {
		return err
	}

	doc, err := client.LoadGOBL(input)
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

func loadClientFromConfig(c *convertOpts) (*fatturapa.Client, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	countryCode := os.Getenv("TRANSMITTER_COUNTRY_CODE")
	taxID := os.Getenv("TRANSMITTER_TAX_ID")

	if countryCode == "" {
		return nil, fmt.Errorf("TRANSMITTER_COUNTRY_CODE not set in .env file")
	}

	if taxID == "" {
		return nil, fmt.Errorf("TRANSMITTER_TAX_ID not set in .env file")
	}

	transmitter := fatturapa.Transmitter{
		CountryCode: countryCode,
		TaxID:       taxID,
	}

	var opts []fatturapa.Option

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

	return fatturapa.NewClient(
		&transmitter,
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
