package encoding

import (
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
)

func EncodeNodeEntry(n entity.Node) (TLV, error) {
	var tlvs []TLV

	// Local ID
	localIDTLV, err := NewTLV(uint8(constants.NODE_ID), EncodeByte(n.ID))
	if err != nil {
		return nil, err
	}
	tlvs = append(tlvs, localIDTLV)

	// Sensor Entries
	panic("not implemented")
}
