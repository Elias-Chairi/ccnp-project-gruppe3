package messages

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/selectors"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

// commandMessage represents a COMMAND message (Type 0x30)
type commandMessage struct {
	NodeSelector     *selectors.NodeSelector
	ActuatorSelector selectors.ActuatorSelector
	ActuatorState    any
}

// NewCommandMessage creates a new COMMAND message
func NewCommandMessage(actuatorSelector selectors.ActuatorSelector, actuatorState any) *commandMessage {
	return &commandMessage{
		NodeSelector:     nil,
		ActuatorSelector: actuatorSelector,
		ActuatorState:    actuatorState,
	}
}

// NewCommandMessageWithNode creates a new COMMAND message with a Node Selector
func NewCommandMessageWithNode(nodeSelector selectors.NodeSelector, actuatorSelector selectors.ActuatorSelector, actuatorState any) *commandMessage {
	return &commandMessage{
		NodeSelector:     &nodeSelector,
		ActuatorSelector: actuatorSelector,
		ActuatorState:    actuatorState,
	}
}

// Encode encodes the COMMAND message to bytes
func (m *commandMessage) Encode() ([]byte, error) {
	var tlvs []tlv.TLV

	// Encode node selector
	if m.NodeSelector != nil {
		nodeSelectorTLV, err := m.NodeSelector.Encode()
		if err != nil {
			return nil, fmt.Errorf("failed to encode Node Selector: %w", err)
		}
		tlvs = append(tlvs, nodeSelectorTLV)
	}

	// Encode actuator selector
	actuatorSelectorTLV, err := m.ActuatorSelector.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode Actuator Selector: %w", err)
	}
	tlvs = append(tlvs, actuatorSelectorTLV)

	// Encode actuator state
	actuatorStateValueTLV, err := protocol.EncodeAny(m.ActuatorState)
	if err != nil {
		return nil, fmt.Errorf("failed to encode Actuator State: %w", err)
	}
	actuatorStateTLV, err := tlv.NewTLV(uint8(constants.ACTUATOR_STATE), actuatorStateValueTLV.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to create Actuator State TLV: %w", err)
	}
	tlvs = append(tlvs, actuatorStateTLV)

	value := tlv.EncodeMultipleTLVs(tlvs)
	mainTLV, err := tlv.NewTLV(uint8(constants.COMMAND), value)
	if err != nil {
		return nil, fmt.Errorf("failed to create COMMAND TLV: %w", err)
	}
	return mainTLV.Encode(), nil
}

// GetType returns the message type
func (m *commandMessage) Type() constants.MessageType {
	return constants.COMMAND
}

// DecodeCommandMessage decodes a COMMAND message TLV into a commandMessage.
//
// Expected inner TLVs:
//   - Node Selector TLV (required if expectNode is true)
//   - Actuator Selector TLV (required)
//   - Actuator State TLV (required)
func DecodeCommandMessage(t tlv.TLV, expectNode bool) (commandMessage, error) {
	if t == nil {
		return commandMessage{}, fmt.Errorf("TLV is nil")
	}
	if t.Type() != uint8(constants.COMMAND) {
		return commandMessage{}, fmt.Errorf("expected COMMAND type, got 0x%x", t.Type())
	}

	nested, err := tlv.DecodeMultipleTLVs(t.Value())
	if err != nil {
		return commandMessage{}, fmt.Errorf("failed to decode nested TLVs: %w", err)
	}

	var (
		nodeSel  *selectors.NodeSelector
		actSel   *selectors.ActuatorSelector
		actState any
	)

	for _, inner := range nested {
		switch {
		case constants.NodeSelector(inner.Type()).IsValid():
			ns, err := selectors.DecodeNodeSelector(inner)
			if err != nil {
				return commandMessage{}, fmt.Errorf("failed to decode Node Selector: %w", err)
			}
			nodeSel = &ns

		case constants.ActuatorSelector(inner.Type()).IsValid():
			as, err := selectors.DecodeActuatorSelector(inner)
			if err != nil {
				return commandMessage{}, fmt.Errorf("failed to decode Actuator Selector: %w", err)
			}
			actSel = &as

		case inner.Type() == uint8(constants.ACTUATOR_STATE):
			state, err := protocol.DecodeAny(inner)
			if err != nil {
				return commandMessage{}, fmt.Errorf("failed to decode Actuator State: %w", err)
			}
			actState = state

		default:
			return commandMessage{}, fmt.Errorf("unexpected TLV type 0x%x in COMMAND message", inner.Type())
		}
	}

	// --- Validation ---
	if expectNode && nodeSel == nil {
		return commandMessage{}, fmt.Errorf("missing required Node Selector TLV")
	}
	if actSel == nil {
		return commandMessage{}, fmt.Errorf("missing required Actuator Selector TLV")
	}
	if actState == nil {
		return commandMessage{}, fmt.Errorf("missing required Actuator State TLV")
	}

	return commandMessage{nodeSel, *actSel, actState}, nil
}
