package main

import (
	"fmt"

	fatturapa "github.com/invopop/gobl.fatturapa"
	"github.com/spf13/cobra"
)

type convertOpts struct {
	*rootOpts
	cert     string
	password string
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

	// TODO: where should the client be injected when using the CLI?
	client := fatturapa.NewClient("IT", "01234567890")

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
