// Package main implements the CLI as well as mage commands (toplevel mage.go)
package main

import (
	"encoding/json"
	"fmt"
	"os"

	fatturapa "github.com/invopop/gobl.fatturapa"
	"github.com/spf13/cobra"
)

type toGOBLOpts struct {
	*convertOpts
	// Add any specific options for XML-to-GOBL conversion here
	pretty bool
}

func toGOBL(c *convertOpts) *toGOBLOpts {
	return &toGOBLOpts{convertOpts: c}
}

func (t *toGOBLOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "to-gobl [infile] [outfile]",
		Short: "Convert a FatturaPA XML into a GOBL JSON",
		Long:  `Convert from FatturaPA XML format to GOBL JSON format.`,
		RunE:  t.runE,
	}

	// Add any specific flags for XML-to-GOBL conversion
	cmd.Flags().BoolVarP(&t.pretty, "pretty", "p", true, "Output pretty-printed JSON")

	return cmd
}

func (t *toGOBLOpts) runE(_ *cobra.Command, args []string) error {
	input := inputFilename(args)

	data, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	// Convert from XML to GOBL
	env, err := fatturapa.Parse(data)
	if err != nil {
		return fmt.Errorf("converting from FatturaPA to GOBL: %w", err)
	}

	// Marshal to JSON
	if t.pretty {
		data, err = json.MarshalIndent(env, "", "  ")
	} else {
		data, err = json.Marshal(env)
	}
	if err != nil {
		return fmt.Errorf("marshaling to JSON: %w", err)
	}

	outFile := outputFilename(args)
	// Write the output
	if err = os.WriteFile(outFile, data, 0644); err != nil {
		return fmt.Errorf("writing JSON output: %w", err)
	}

	return nil
}
