// Package main implements the CLI as well as mage commands (toplevel mage.go)
package main

import (
	"github.com/spf13/cobra"
)

type convertOpts struct {
	*rootOpts
}

func convert(o *rootOpts) *convertOpts {
	return &convertOpts{rootOpts: o}
}

func (c *convertOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert between GOBL and FatturaPA formats",
		Long:  `Convert between GOBL JSON and FatturaPA XML formats in both directions.`,
	}

	// Add subcommands
	cmd.AddCommand(toXML(c).cmd())
	cmd.AddCommand(fromXML(c).cmd())

	return cmd
}
