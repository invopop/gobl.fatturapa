package fatturapa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPostOfficeBox(t *testing.T) {
	tests := []struct {
		name   string
		street string
		want   bool
	}{
		{"P.O. prefix", "P.O. 123", true},
		{"PO Box prefix", "PO Box 123", true},
		{"P.O.Box prefix", "P.O.Box 123", true},
		{"P.O Box prefix", "P.O Box 123", true},
		{"PO BOX prefix", "PO BOX 123", true},
		{"empty", "", false},
		{"len 5 non-match", "BOX12", false},
		{"len 6 non-match", "BOX123", false},
		{"len 7 non-match", "ROMA 1 ", false},
		{"PO prefix inside string", "Apt 3, PO Box 12", false},
		{"lowercase p.o. not matched", "p.o. 123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isPostOfficeBox(tt.street))
		})
	}
}
