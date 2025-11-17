package messages

import (
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
)

type nodeRemoveMessage struct {
	NodeID uint8
}

func NewNodeRemovedMessage(nodeID uint8) *nodeRemoveMessage {
	return &nodeRemoveMessage{
		NodeID: nodeID,
	}
}

func (m *nodeRemoveMessage) Type() constants.MessageType {
	return constants.NODE_REMOVED
}

func (m *nodeRemoveMessage) Encode() (encoding.TLV, error) {
	nodeID, err := encoding.NewTLV(uint8(constants.SINGLE_NODE), encoding.EncodeByte(m.NodeID))
	if err != nil {
		return nil, err
	}

	mainTLV, err := encoding.NewTLV(uint8(constants.NODE_REMOVED), nodeID.Encode())
	if err != nil {
		return nil, err
	}

	return mainTLV, nil
}
