package encoding

import (
	"encoding/binary"
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
