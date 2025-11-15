package messages

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/selectors"
)

type ActuatorUpdateMessage struct {
	ActuatorSelector selectors.ActuatorSelector
	ActuatorState    any
}

func NewActuatorUpdateMessage(actuatorSelector selectors.ActuatorSelector, actuatorState any) *ActuatorUpdateMessage {
	return &ActuatorUpdateMessage{
		ActuatorSelector: actuatorSelector,
		ActuatorState:    actuatorState,
	}
}

// Type returns the message type.
func (m *ActuatorUpdateMessage) Type() constants.MessageType {
	return constants.ACTUATOR_UPDATE
}

func (m *ActuatorUpdateMessage) Encode() (encoding.TLV, error) {
	var tlvs []encoding.TLV

	// Encode actuator selector
	actuatorSelectorTLV, err := m.ActuatorSelector.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode Actuator Selector: %w", err)
	}
	tlvs = append(tlvs, actuatorSelectorTLV)

	// Encode actuator state
	actuatorStateValueTLV, err := encoding.EncodeAny(m.ActuatorState)
	if err != nil {
		return nil, fmt.Errorf("failed to encode Actuator State: %w", err)
	}
	actuatorStateTLV, err := encoding.NewTLV(uint8(constants.ACTUATOR_STATE), actuatorStateValueTLV.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to create Actuator State TLV: %w", err)
	}
	tlvs = append(tlvs, actuatorStateTLV)

	value := encoding.EncodeMultipleTLVs(tlvs)
	mainTLV, err := encoding.NewTLV(uint8(constants.COMMAND), value)
	if err != nil {
		return nil, fmt.Errorf("failed to create COMMAND TLV: %w", err)
	}
	return mainTLV, nil
}

func DecodeActuatorUpdateMessage(t encoding.TLV) (*ActuatorUpdateMessage, error) {
	if t == nil {
		return nil, fmt.Errorf("TLV is nil")
	}
	if t.Type() != uint8(constants.ACTUATOR_UPDATE) {
		return nil, fmt.Errorf("expected ACTUATOR_UPDATE type, got 0x%x", t.Type())
	}

	tlvs, err := encoding.DecodeMultipleTLVs(t.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to decode nested TLVs: %w", err)
	}

	var (
		actSel   *selectors.ActuatorSelector
		actState any
	)

	for _, inner := range tlvs {
		switch {
		case constants.ActuatorSelector(inner.Type()).IsValid():
			if actSel != nil {
				return nil, fmt.Errorf("duplicate Actuator Selector TLV found")
			}
			as, err := selectors.DecodeActuatorSelector(inner)
			if err != nil {
				return nil, fmt.Errorf("failed to decode Actuator Selector: %w", err)
			}
			actSel = &as

		case inner.Type() == uint8(constants.ACTUATOR_STATE):
			if actState != nil {
				return nil, fmt.Errorf("duplicate Actuator State TLV found")
			}
			innerTLV, err := encoding.DecodeTLV(inner.Value())
			if err != nil {
				return nil, fmt.Errorf("failed to decode Actuator State inner TLV: %w", err)
			}
			actState, err = encoding.DecodeAny(innerTLV)
			if err != nil {
				return nil, fmt.Errorf("failed to decode Actuator State: %w", err)
			}

		default:
			return nil, fmt.Errorf("unexpected TLV type 0x%x in COMMAND message", inner.Type())
		}
	}

	if actSel == nil {
		return nil, fmt.Errorf("missing required Actuator Selector TLV")
	}
	if actState == nil {
		return nil, fmt.Errorf("missing required Actuator State TLV")
	}

	return &ActuatorUpdateMessage{*actSel, actState}, nil
}
