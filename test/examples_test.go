package test_test

import (
	"testing"

	"github.com/invopop/gobl.fatturapa/test"
	"github.com/stretchr/testify/assert"
)

func TestExamples(t *testing.T) {
	err := test.TestConversion()
	assert.NoError(t, err)
}
