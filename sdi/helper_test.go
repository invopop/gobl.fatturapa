package sdi_test

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	sdi "github.com/invopop/gobl.fatturapa/sdi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func handlerFunc(env *sdi.Envelope) {
	if env.Body.FileSubmissionMetadata != nil {
		log.Printf("parsing MetadatiInvioFile:\n")
	}
	if env.Body.NonDeliveryNotificationMessage != nil {
		log.Printf("parsing NotificaMancataConsegna:\n")
	}
	if env.Body.InvoiceTransmissionCertificate != nil {
		log.Printf("parsing AttestazioneTrasmissioneFattura:\n")
	}
}

func TestParseMessage(t *testing.T) {
	t.Run("parse MetadatiInvioFile", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()

		reader, err := os.Open("./test/examples/ESB85905495_00010_MT_001.xml")
		require.NoError(t, err)

		err = sdi.ParseMessage(io.NopCloser(reader), handlerFunc)
		require.NoError(t, err)

		assert.Contains(t, buf.String(), "parsing MetadatiInvioFile")
	})

	t.Run("parse NotificaMancataConsegna", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()

		reader, err := os.Open("./test/examples/ESB85905495_00010_MT_004.xml")
		require.NoError(t, err)

		err = sdi.ParseMessage(io.NopCloser(reader), handlerFunc)
		require.NoError(t, err)

		assert.Contains(t, buf.String(), "parsing NotificaMancataConsegna")
	})

	t.Run("parse AttestazioneTrasmissioneFattura", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()

		reader, err := os.Open("./test/examples/ESB85905495_00010_MT_005.xml")
		require.NoError(t, err)

		err = sdi.ParseMessage(io.NopCloser(reader), handlerFunc)
		require.NoError(t, err)

		assert.Contains(t, buf.String(), "parsing AttestazioneTrasmissioneFattura")
	})
}
