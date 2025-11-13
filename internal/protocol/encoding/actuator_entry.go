package encoding

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
)

// EncodeActuatorEntry encodes an actuator into an ACTUATOR_ENTRY TLV with nested fields.
func EncodeActuatorEntry(a entity.Actuator[any]) (TLV, error) {
	var tlvs []TLV

	// Actuator ID
	aid, err := NewTLV(uint8(constants.ACTUATOR_ID), EncodeByte(a.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to create ACTUATOR_ID TLV: %w", err)
	}
	tlvs = append(tlvs, aid)

	// Actuator Type
	at, err := NewTLV(uint8(constants.ACTUATOR_TYPE_FIELD), []byte(a.Type))
	if err != nil {
		return nil, fmt.Errorf("failed to create ACTUATOR_TYPE TLV: %w", err)
	}
	tlvs = append(tlvs, at)

	// Actuator Unit
	au, err := NewTLV(uint8(constants.ACTUATOR_UNIT), []byte(a.Unit))
	if err != nil {
		return nil, fmt.Errorf("failed to create ACTUATOR_UNIT TLV: %w", err)
	}
	tlvs = append(tlvs, au)

	// Actuator State (with data type TLV)
	stateTLV, err := EncodeAny(a.State)
	if err != nil {
		return nil, fmt.Errorf("failed to encode actuator state: %w", err)
	}
	as, err := NewTLV(uint8(constants.ACTUATOR_STATE), stateTLV.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to create ACTUATOR_STATE TLV: %w", err)
	}
	tlvs = append(tlvs, as)

	// Wrap all fields in ACTUATOR_ENTRY TLV
	tlv, err := NewTLV(uint8(constants.ACTUATOR_ENTRY), EncodeMultipleTLVs(tlvs))
	if err != nil {
		return nil, fmt.Errorf("failed to create ACTUATOR_ENTRY TLV: %w", err)
	}
	return tlv, nil
}

// DecodeActuatorEntry decodes an ACTUATOR_ENTRY TLV into an Actuator and validates field lengths.
func DecodeActuatorEntry(t TLV) (*entity.Actuator[any], error) {
	if t == nil {
		return nil, fmt.Errorf("data is nil")
	}

	if t.Type() != uint8(constants.ACTUATOR_ENTRY) {
		return nil, fmt.Errorf("expected ACTUATOR_ENTRY type, got %x", t.Type())
	}

	tlvs, err := DecodeMultipleTLVs(t.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to decode multiple TLVs: %w", err)
	}

	a := &entity.Actuator[any]{}
	for _, innerTLV := range tlvs {
		switch constants.ActuatorField(innerTLV.Type()) {
		case constants.ACTUATOR_ID:
			id, err := DecodeByte(innerTLV.Value())
			if err != nil {
				return nil, fmt.Errorf("failed to decode actuator ID: %w", err)
			}
			a.ID = id
		case constants.ACTUATOR_TYPE_FIELD:
			a.Type = string(innerTLV.Value())
		case constants.ACTUATOR_UNIT:
			a.Unit = string(innerTLV.Value())
		case constants.ACTUATOR_STATE:
			stateTLV, err := DecodeTLV(innerTLV.Value())
			if err != nil {
				return nil, fmt.Errorf("failed to decode actuator state TLV: %w", err)
			}
			a.State, err = DecodeAny(stateTLV)
			if err != nil {
				return nil, fmt.Errorf("failed to decode actuator state: %w", err)
			}
		}
	}

	return a, nil
}
