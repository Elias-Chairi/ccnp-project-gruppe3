package selectors

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

// ActuatorSelector represents different ways to select actuators
type ActuatorSelector struct {
	Type        constants.ActuatorSelector // SINGLE_ACTUATOR, ACTUATOR_LIST, ACTUATOR_TYPE, or ALL_ACTUATORS
	ActuatorIDs []uint8                    // For SINGLE_ACTUATOR and ACTUATOR_LIST
	TypeName    string                     // For ACTUATOR_TYPE
}

// NewSingleActuatorSelector creates a selector for a single actuator
func NewSingleActuatorSelector(actuatorID uint8) *ActuatorSelector {
	return &ActuatorSelector{
		Type:        constants.SINGLE_ACTUATOR,
		ActuatorIDs: []uint8{actuatorID},
	}
}

// NewActuatorListSelector creates a selector for multiple actuators
func NewActuatorListSelector(actuatorIDs []uint8) *ActuatorSelector {
	return &ActuatorSelector{
		Type:        constants.ACTUATOR_LIST,
		ActuatorIDs: actuatorIDs,
	}
}

// NewActuatorTypeSelector creates a selector for actuators by type
func NewActuatorTypeSelector(typeName string) *ActuatorSelector {
	return &ActuatorSelector{
		Type:     constants.ACTUATOR_TYPE,
		TypeName: typeName,
	}
}

// NewAllActuatorsSelector creates a selector for all actuators
func NewAllActuatorsSelector() *ActuatorSelector {
	return &ActuatorSelector{
		Type: constants.ALL_ACTUATORS,
	}
}

// Encode encodes the actuator selector as TLV
func (a *ActuatorSelector) Encode() (tlv.TLV, error) {
	switch a.Type {
	case constants.SINGLE_ACTUATOR:
		if len(a.ActuatorIDs) != 1 {
			return nil, fmt.Errorf("SINGLE_ACTUATOR selector must have exactly one actuator ID")
		}
		return tlv.NewTLV(uint8(constants.SINGLE_ACTUATOR), encoding.EncodeByte(a.ActuatorIDs[0]))
	case constants.ACTUATOR_LIST:
		if len(a.ActuatorIDs) == 0 {
			return nil, fmt.Errorf("ACTUATOR_LIST selector must have at least one actuator ID")
		}
		return tlv.NewTLV(uint8(constants.ACTUATOR_LIST), encoding.EncodeByteList(a.ActuatorIDs))
	case constants.ACTUATOR_TYPE:
		if a.TypeName == "" {
			return nil, fmt.Errorf("ACTUATOR_TYPE selector must have a type name")
		}
		return tlv.NewTLV(uint8(constants.ACTUATOR_TYPE), []byte(a.TypeName))
	case constants.ALL_ACTUATORS:
		return tlv.NewTLV(uint8(constants.ALL_ACTUATORS), []byte{})
	default:
		return nil, fmt.Errorf("unknown actuator selector type: %x", a.Type)
	}
}

// DecodeActuatorSelector decodes an actuator selector from TLV
func DecodeActuatorSelector(tlv tlv.TLV) (ActuatorSelector, error) {
	switch constants.ActuatorSelector(tlv.Type()) {
	case constants.SINGLE_ACTUATOR:
		if tlv.Length() != 1 {
			return ActuatorSelector{}, fmt.Errorf("SINGLE_ACTUATOR selector must have exactly one byte")
		}
		return ActuatorSelector{
			Type:        constants.SINGLE_ACTUATOR,
			ActuatorIDs: []uint8{tlv.Value()[0]},
		}, nil
	case constants.ACTUATOR_LIST:
		if tlv.Length() == 0 {
			return ActuatorSelector{}, fmt.Errorf("ACTUATOR_LIST selector must have at least one byte")
		}
		return ActuatorSelector{
			Type:        constants.ACTUATOR_LIST,
			ActuatorIDs: encoding.DecodeByteList(tlv.Value()),
		}, nil
	case constants.ACTUATOR_TYPE:
		if tlv.Length() == 0 {
			return ActuatorSelector{}, fmt.Errorf("ACTUATOR_TYPE selector must have a type name")
		}
		return ActuatorSelector{
			Type:     constants.ACTUATOR_TYPE,
			TypeName: string(tlv.Value()),
		}, nil
	case constants.ALL_ACTUATORS:
		if tlv.Length() != 0 {
			return ActuatorSelector{}, fmt.Errorf("ALL_ACTUATORS selector must have empty value")
		}
		return ActuatorSelector{
			Type: constants.ALL_ACTUATORS,
		}, nil
	default:
		return ActuatorSelector{}, fmt.Errorf("unknown actuator selector type: %x", tlv.Type())
	}
}
