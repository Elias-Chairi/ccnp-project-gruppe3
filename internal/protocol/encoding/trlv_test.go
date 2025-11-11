package encoding_test

import (
	"testing"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/stretchr/testify/assert"
)

// -------------------------------------- Positive tests --------------------------------------

func TestEncodeTRLV(t *testing.T) {
	tl, _ := encoding.NewTLV(0x10, []byte{0xDE, 0xAD, 0xBE, 0xEF})
	trlv := encoding.TRLV{
		TLV:       tl,
		RequestID: 0x1234,
	}

	expected := []byte{0x10, 0x12, 0x34, 0x00, 0x04, 0xDE, 0xAD, 0xBE, 0xEF}
	assert.Equal(t, expected, trlv.Encode())
}
