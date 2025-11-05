package encoding

import (
	"encoding/binary"
	"fmt"
	"io"
)

// TRLV represents a Type-RequestID-Length-Value structure.
type TRLV struct {
	TLV       TLV
	RequestID uint16
}

// Encode encodes the TRLV to bytes.
//
// Format: [Type(1)][RequestID(2)][Length(2)][Value(N)]
func (t TRLV) Encode() []byte {
	if t.TLV == nil {
		return nil
	}
	result := make([]byte, 5+len(t.TLV.Value()))
	result[0] = t.TLV.Type()
	binary.BigEndian.PutUint16(result[1:3], t.RequestID)
	binary.BigEndian.PutUint16(result[3:5], t.TLV.Length())
	copy(result[5:], t.TLV.Value())
	return result
}

// ReadTRLV reads a TRLV from the given reader.
//
// It first reads until it has recived the full header (5 bytes),
// then reads until the full value is received based on the length field in the header.
//
// Returns the TRLV, number of bytes read, and an error if any.
func ReadTRLV(r io.Reader) (*TRLV, int, error) {
	header := make([]byte, 5) // expecting Type(1) + RequestID(2) + Length(2)
	n, err := io.ReadFull(r, header)
	if err != nil {
		return nil, n, fmt.Errorf("error reading TRLV header: %w", err)
	}

	length := binary.BigEndian.Uint16(header[3:5])
	value := make([]byte, length)
	n, err = io.ReadFull(r, value) // expecting Value(length)
	if err != nil {
		return nil, n + 5, fmt.Errorf("error reading TRLV value: %w", err)
	}

	return &TRLV{
		TLV: &tlv{
			tlvType: header[0],
			length:  length,
			value:   value,
		},
		RequestID: binary.BigEndian.Uint16(header[1:3]),
	}, n + 5, nil
}
