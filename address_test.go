package fatturapa

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestAddressRegion(t *testing.T) {
	t.Run("should return the region two letter code", func(t *testing.T) {
		addr := &org.Address{
			Number:   "1",
			Street:   "Via Roma",
			Locality: "Roma",
			Region:   "RM",
			Code:     "00100",
			Country:  l10n.IT,
		}

		out := newAddress(addr)
		assert.Equal(t, "RM", out.Provincia)
	})

	t.Run("should ignore text name", func(t *testing.T) {
		addr := &org.Address{
			Number:   "1",
			Street:   "Via Roma",
			Locality: "Roma",
			Region:   "Rome",
			Code:     "00100",
			Country:  l10n.IT,
		}

		out := newAddress(addr)
		assert.Empty(t, out.Provincia)
	})

	t.Run("should ignore foreign addresses", func(t *testing.T) {
		addr := &org.Address{
			Number:   "2",
			Street:   "Rome Street",
			Locality: "London",
			Region:   "RM",
			Code:     "00100",
			Country:  l10n.GB,
		}

		out := newAddress(addr)
		assert.Empty(t, out.Provincia)
	})
}
