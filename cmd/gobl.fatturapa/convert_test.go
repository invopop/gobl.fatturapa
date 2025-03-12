package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_convert(t *testing.T) {
	cmd := convert(root()).cmd()
	assert.Equal(t, "convert", cmd.Use)
	assert.Equal(t, 2, len(cmd.Commands()))

	// Check that both subcommands are registered
	var hasToXML, hasFromXML bool
	for _, subcmd := range cmd.Commands() {
		switch subcmd.Use {
		case "to-xml [infile] [outfile]":
			hasToXML = true
		case "from-xml [infile] [outfile]":
			hasFromXML = true
		}
	}
	assert.True(t, hasToXML, "to-xml subcommand should be registered")
	assert.True(t, hasFromXML, "from-xml subcommand should be registered")
}

func Test_toXML(t *testing.T) {
	cmd := toXML(convert(root())).cmd()
	assert.Equal(t, "to-xml [infile] [outfile]", cmd.Use)

	// Check that all flags are registered
	flags := cmd.Flags()
	certFlag := flags.Lookup("cert")
	assert.NotNil(t, certFlag, "cert flag should be registered")

	passwordFlag := flags.Lookup("password")
	assert.NotNil(t, passwordFlag, "password flag should be registered")

	transmitterFlag := flags.Lookup("transmitter")
	assert.NotNil(t, transmitterFlag, "transmitter flag should be registered")

	timestampFlag := flags.Lookup("with-timestamp")
	assert.NotNil(t, timestampFlag, "with-timestamp flag should be registered")
}

func Test_fromXML(t *testing.T) {
	cmd := fromXML(convert(root())).cmd()
	assert.Equal(t, "from-xml [infile] [outfile]", cmd.Use)

	// Check that all flags are registered
	flags := cmd.Flags()
	prettyFlag := flags.Lookup("pretty")
	assert.NotNil(t, prettyFlag, "pretty flag should be registered")
}
