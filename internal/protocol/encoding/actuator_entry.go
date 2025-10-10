package encoding

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

// EncodeActuatorEntry encodes an actuator into an ACTUATOR_ENTRY TLV with nested fields.
func EncodeActuatorEntry(a entity.Actuator[any]) (tlv.TLV, error) {
	var tlvs []tlv.TLV

	// Actuator ID
	aid, err := tlv.NewTLV(uint8(constants.ACTUATOR_ID), EncodeByte(a.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to create ACTUATOR_ID TLV: %w", err)
	}
	tlvs = append(tlvs, aid)

	// Actuator Type
	at, err := tlv.NewTLV(uint8(constants.ACTUATOR_TYPE_FIELD), []byte(a.Type))
	if err != nil {
		return nil, fmt.Errorf("failed to create ACTUATOR_TYPE TLV: %w", err)
	}
	tlvs = append(tlvs, at)

	// Actuator Unit
	au, err := tlv.NewTLV(uint8(constants.ACTUATOR_UNIT), []byte(a.Unit))
	if err != nil {
		return nil, fmt.Errorf("failed to create ACTUATOR_UNIT TLV: %w", err)
	}
	tlvs = append(tlvs, au)

	// Actuator State (with data type TLV)
	stateTLV, err := EncodeAny(a.State)
	if err != nil {
		return nil, fmt.Errorf("failed to encode actuator state: %w", err)
	}
	as, err := tlv.NewTLV(uint8(constants.ACTUATOR_STATE), stateTLV.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to create ACTUATOR_STATE TLV: %w", err)
	}
	tlvs = append(tlvs, as)

	// Wrap all fields in ACTUATOR_ENTRY TLV
	tlv, err := tlv.NewTLV(uint8(constants.ACTUATOR_ENTRY), tlv.EncodeMultipleTLVs(tlvs))
	if err != nil {
		return nil, fmt.Errorf("failed to create ACTUATOR_ENTRY TLV: %w", err)
	}
	return tlv, nil
}

// DecodeActuatorEntry decodes an ACTUATOR_ENTRY TLV into an Actuator and validates field lengths.
func DecodeActuatorEntry(t tlv.TLV) (*entity.Actuator[any], error) {
	if t == nil {
		return nil, fmt.Errorf("data is nil")
	}

	if t.Type() != uint8(constants.ACTUATOR_ENTRY) {
		return nil, fmt.Errorf("expected ACTUATOR_ENTRY type, got %x", t.Type())
	}

	tlvs, err := tlv.DecodeMultipleTLVs(t.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to decode multiple TLVs: %w", err)
	}

	a := &entity.Actuator[any]{}
	for _, tlv := range tlvs {
		switch constants.ActuatorField(tlv.Type()) {
		case constants.ACTUATOR_ID:
			if tlv.Length() != 1 {
				return nil, fmt.Errorf("invalid actuator ID length")
			}
			a.ID = tlv.Value()[0]
		case constants.ACTUATOR_TYPE_FIELD:
			a.Type = string(tlv.Value())
		case constants.ACTUATOR_UNIT:
			a.Unit = string(tlv.Value())
		case constants.ACTUATOR_STATE:
			if tlv.Length() != 1 {
				return nil, fmt.Errorf("invalid actuator state length")
			}
			a.State = tlv.Value()[0]
		}
	}

	return a, nil
}
