package selectors

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

// NodeSelector represents different ways to select nodes
type NodeSelector struct {
	Type    constants.NodeSelector // SINGLE_NODE, NODE_LIST, or ALL_NODES
	NodeIDs []uint8                // For SINGLE_NODE and NODE_LIST
}

// NewSingleNodeSelector creates a selector for a single node
func NewSingleNodeSelector(nodeID uint8) *NodeSelector {
	return &NodeSelector{
		Type:    constants.SINGLE_NODE,
		NodeIDs: []uint8{nodeID},
	}
}

// NewNodeListSelector creates a selector for multiple nodes
func NewNodeListSelector(nodeIDs []uint8) *NodeSelector {
	return &NodeSelector{
		Type:    constants.NODE_LIST,
		NodeIDs: nodeIDs,
	}
}

// NewAllNodesSelector creates a selector for all nodes
func NewAllNodesSelector() *NodeSelector {
	return &NodeSelector{
		Type:    constants.ALL_NODES,
		NodeIDs: nil,
	}
}

// Encode encodes the node selector as TLV
func (n *NodeSelector) Encode() (tlv.TLV, error) {
	switch n.Type {
	case constants.SINGLE_NODE:
		if len(n.NodeIDs) != 1 {
			return nil, fmt.Errorf("SINGLE_NODE selector must have exactly one node ID")
		}
		return tlv.NewTLV(uint8(constants.SINGLE_NODE), protocol.EncodeByte(n.NodeIDs[0]))
	case constants.NODE_LIST:
		if len(n.NodeIDs) == 0 {
			return nil, fmt.Errorf("NODE_LIST selector must have at least one node ID")
		}
		return tlv.NewTLV(uint8(constants.NODE_LIST), protocol.EncodeByteList(n.NodeIDs))
	case constants.ALL_NODES:
		return tlv.NewTLV(uint8(constants.ALL_NODES), []byte{})
	default:
		return nil, fmt.Errorf("unknown node selector type: %x", n.Type)
	}
}

// DecodeNodeSelector decodes a node selector from TLV
func DecodeNodeSelector(tlv tlv.TLV) (NodeSelector, error) {
	switch constants.NodeSelector(tlv.Type()) {
	case constants.SINGLE_NODE:
		if tlv.Length() != 1 {
			return NodeSelector{}, fmt.Errorf("SINGLE_NODE selector must have exactly one byte")
		}
		return NodeSelector{
			Type:    constants.SINGLE_NODE,
			NodeIDs: []uint8{tlv.Value()[0]},
		}, nil
	case constants.NODE_LIST:
		if tlv.Length() == 0 {
			return NodeSelector{}, fmt.Errorf("NODE_LIST selector must have at least one byte")
		}
		return NodeSelector{
			Type:    constants.NODE_LIST,
			NodeIDs: protocol.DecodeByteList(tlv.Value()),
		}, nil
	case constants.ALL_NODES:
		if tlv.Length() != 0 {
			return NodeSelector{}, fmt.Errorf("ALL_NODES selector must have empty value")
		}
		return NodeSelector{
			Type:    constants.ALL_NODES,
			NodeIDs: nil,
		}, nil
	default:
		return NodeSelector{}, fmt.Errorf("unknown node selector type: %x", tlv.Type())
	}
}
