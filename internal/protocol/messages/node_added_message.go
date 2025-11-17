package messages

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
)

type nodeAddedMessage struct {
	Node entity.Node
}

func NewNodeAddedMessage(node entity.Node) *nodeAddedMessage {
	return &nodeAddedMessage{
		Node: node,
	}
}

func (m *nodeAddedMessage) Type() constants.MessageType {
	return constants.NODE_ADDED
}

func (m *nodeAddedMessage) Encode() (encoding.TLV, error) {
	nodeEntry, err := encoding.EncodeNodeEntry(m.Node)
	if err != nil {
		return nil, err
	}

	mainTLV, err := encoding.NewTLV(uint8(constants.NODE_ADDED), nodeEntry.Encode())
	if err != nil {
		return nil, err
	}

	return mainTLV, nil
}

func DecodeNodeAddedMessage(tlv encoding.TLV) (*nodeAddedMessage, error) {
	if tlv.Type() != uint8(constants.NODE_ADDED) {
		return nil, fmt.Errorf("expected NODE_ADDED TLV, got %d", tlv.Type())
	}

	nodeAddedVal, err := encoding.DecodeTLV(tlv.Value())
	if err != nil {
		return nil, err
	}

	node, err := encoding.DecodeNodeEntry(nodeAddedVal)
	if err != nil {
		return nil, err
	}
	return &nodeAddedMessage{
		Node: *node,
	}, nil
}
