// Package main implements the CLI as well as mage commands (toplevel mage.go)
package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	fatturapa "github.com/invopop/gobl.fatturapa"
	"github.com/spf13/cobra"
)

type fromXMLOpts struct {
	*convertOpts
	// Add any specific options for XML-to-GOBL conversion here
	prettyJSON bool
}

func fromXML(c *convertOpts) *fromXMLOpts {
	return &fromXMLOpts{convertOpts: c}
}

func (f *fromXMLOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "from-xml [infile] [outfile]",
		Short: "Convert a FatturaPA XML into a GOBL JSON",
		Long:  `Convert from FatturaPA XML format to GOBL JSON format.`,
		RunE:  f.runE,
	}

	// Add any specific flags for XML-to-GOBL conversion
	cmd.Flags().BoolVarP(&f.prettyJSON, "pretty", "p", true, "Output pretty-printed JSON")

	return cmd
}

func (f *fromXMLOpts) runE(cmd *cobra.Command, args []string) error {
	input, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer input.Close() // nolint:errcheck

	out, err := f.openOutput(cmd, args)
	if err != nil {
		return err
	}
	defer out.Close() // nolint:errcheck

	// Create a new converter
	converter := fatturapa.NewConverter()

	// Read the XML input
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(input); err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	// Convert from XML to GOBL
	env, err := converter.ConvertToGOBL(buf.Bytes())
	if err != nil {
		return fmt.Errorf("converting from FatturaPA to GOBL: %w", err)
	}

	// Marshal to JSON
	var data []byte
	if f.prettyJSON {
		data, err = json.MarshalIndent(env, "", "  ")
	} else {
		data, err = json.Marshal(env)
	}
	if err != nil {
		return fmt.Errorf("marshaling to JSON: %w", err)
	}

	// Write the output
	if _, err = out.Write(data); err != nil {
		return fmt.Errorf("writing JSON output: %w", err)
	}

	return nil
}
