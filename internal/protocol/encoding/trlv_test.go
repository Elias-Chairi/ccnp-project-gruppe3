package encoding_test

import (
	"bytes"
	"io"
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

func TestReadTRLV(t *testing.T) {
	data := []byte{0x20, 0x56, 0x78, 0x00, 0x03, 0xCA, 0xFE, 0xBA}
	r := bytes.NewReader(data)
	trlv, n, err := encoding.ReadTRLV(r)
	assert.NoError(t, err)
	assert.Equal(t, 8, n)
	assert.Equal(t, uint8(0x20), trlv.TLV.Type())
	assert.Equal(t, uint16(3), trlv.TLV.Length())
	assert.Equal(t, uint16(0x5678), trlv.RequestID)
	assert.Equal(t, []byte{0xCA, 0xFE, 0xBA}, trlv.TLV.Value())
}

// -------------------------------------- Negative tests --------------------------------------

func TestReadTRLV_IncompleteHeader(t *testing.T) {
	data := []byte{0x30, 0x00} // Incomplete header
	r := bytes.NewReader(data)
	trlv, n, err := encoding.ReadTRLV(r)
	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
	assert.Nil(t, trlv)
	assert.Equal(t, 2, n)
}

func TestReadTRLV_IncompleteValue(t *testing.T) {
	data := []byte{0x40, 0x9A, 0xBC, 0x00, 0x05, 0x11, 0x22} // Incomplete value
	r := bytes.NewReader(data)
	trlv, n, err := encoding.ReadTRLV(r)
	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
	assert.Nil(t, trlv)
	assert.Equal(t, 7, n)
}
