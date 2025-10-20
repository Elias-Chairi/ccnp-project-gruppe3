package tlv_test

import (
	"testing"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
	"github.com/stretchr/testify/assert"
)

// -------------------------------------- Positive tests --------------------------------------

func TestNewTLV(t *testing.T) {
	assert := assert.New(t)

	val := []byte{0xAA, 0xBB}
	tt, err := tlv.NewTLV(0x44, val)
	assert.NoError(err)
	assert.NotNil(tt)

	assert.Equal(uint8(0x44), tt.Type())
	assert.Equal(uint16(2), tt.Length())
	assert.Equal(val, tt.Value())

	// NewTLV stores a reference to the provided slice (no copy)
	val[0] = 0xCC
	assert.Equal(uint8(0xCC), tt.Value()[0])
}

func TestEncode(t *testing.T) {
	val := []byte{0x01, 0x02, 0x03}
	tt, err := tlv.NewTLV(0x12, val)
	assert.NoError(t, err)
	encoded := tt.Encode()

	// Expect [Type][Length][Value] = [0x12][0x00 0x03][0x01 0x02 0x03]
	expected := []byte{0x12, 0x00, 0x03, 0x01, 0x02, 0x03}
	assert.Equal(t, expected, encoded)
}

func TestDecodeTLV(t *testing.T) {
	assert := assert.New(t)

	data := []byte{0x34, 0x00, 0x02, 0x10, 0x20}
	got, err := tlv.DecodeTLV(data)
	assert.NoError(err)
	assert.Equal(uint8(0x34), got.Type())
	assert.Equal(uint16(2), got.Length())
	assert.Equal([]byte{0x10, 0x20}, got.Value())

	// DecodeTLV makes a defensive copy; mutating original buffer shouldn't affect Value
	data[3] = 0xFF
	assert.Equal([]byte{0x10, 0x20}, got.Value())
}

func TestEncodeMultipleTLVs(t *testing.T) {
	a, _ := tlv.NewTLV(0x10, []byte{0xAA})
	b, _ := tlv.NewTLV(0x20, []byte{0xBB, 0xCC})
	expected := append(a.Encode(), b.Encode()...)

	got := tlv.EncodeMultipleTLVs([]tlv.TLV{a, b})
	assert.Equal(t, expected, got)
}

func TestDecodeMultipleTLVs(t *testing.T) {
	assert := assert.New(t)

	a, _ := tlv.NewTLV(0x01, []byte{0x11})
	b, _ := tlv.NewTLV(0x02, []byte{0x22, 0x33})
	concat := append(a.Encode(), b.Encode()...)

	tlvs, err := tlv.DecodeMultipleTLVs(concat)
	assert.NoError(err)
	assert.Len(tlvs, 2)

	assert.Equal(uint8(0x01), tlvs[0].Type())
	assert.Equal([]byte{0x11}, tlvs[0].Value())
	assert.Equal(uint8(0x02), tlvs[1].Type())
	assert.Equal([]byte{0x22, 0x33}, tlvs[1].Value())
}

// -------------------------------------- Negative tests --------------------------------------

func TestNewTLV_EmptyValue(t *testing.T) {
	tt, err := tlv.NewTLV(0x01, nil)
	assert.NoError(t, err)
	assert.NotNil(t, tt)

	tt, err = tlv.NewTLV(0x01, []byte{})
	assert.NoError(t, err)
	assert.NotNil(t, tt)
}

func TestNewTLV_TooLarge(t *testing.T) {
	val := make([]byte, 65536)
	tt, err := tlv.NewTLV(0x01, val)
	assert.Error(t, err)
	assert.Nil(t, tt)
}

func TestDecodeTLV_InvalidInput(t *testing.T) {
	// Too short for header
	_, err := tlv.DecodeTLV([]byte{0x01, 0x00})
	assert.Error(t, err)

	// Declared value longer than available
	_, err = tlv.DecodeTLV([]byte{0x01, 0x00, 0x03, 0xAA})
	assert.Error(t, err)
}

func TestDecodeMultipleTLVs_TrailingBytesIgnored(t *testing.T) {
	a, _ := tlv.NewTLV(0x01, []byte{0x11})
	b, _ := tlv.NewTLV(0x02, []byte{0x22})
	concat := append(a.Encode(), b.Encode()...)
	// Add 2 trailing bytes (< header size) which should be ignored
	concat = append(concat, 0xFF, 0xFF)

	tlvs, err := tlv.DecodeMultipleTLVs(concat)
	assert.NoError(t, err)
	assert.Len(t, tlvs, 2)
}

func TestDecodeMultipleTLVs_ErrorOnSecond(t *testing.T) {
	a, _ := tlv.NewTLV(0x01, []byte{0x11})
	// Malformed second TLV: header says length=2 but only 1 value byte present
	malformed := []byte{0x02, 0x00, 0x02, 0x22}
	concat := append(a.Encode(), malformed...)

	tlvs, err := tlv.DecodeMultipleTLVs(concat)
	assert.Error(t, err)
	assert.Nil(t, tlvs)
}
