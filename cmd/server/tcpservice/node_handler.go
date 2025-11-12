package tcpservice

import (
	"net"
	"sync"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util"
)

var handleNode messageHandler = func(msg *encoding.Message) messages.AckErrorMessage {
	switch msg.TLV.Type() {
	case uint8(constants.SENSOR_UPDATE):
		// todo: handle sensor update
		return messages.AckSuccessMessage()
	case uint8(constants.ACK_ERROR):
		// todo: handle ack error
		return messages.AckSuccessMessage()
	default:
		return messages.AckErrorMessage{
			Code: constants.ERR_INVALID_MESSAGE_TYPE,
		}
	}
}

// NodeRegistry manages node IDs and their associated connections.
type NodeRegistry struct {
	mu    sync.RWMutex
	nodes map[uint8]net.Conn
}

// NewNodeRegistry initializes and returns a new NodeRegistry.
func NewNodeRegistry() *NodeRegistry {
	return &NodeRegistry{
		nodes: make(map[uint8]net.Conn),
	}
}

// CreateNodeID assigns a unique node ID and stores the connection.
func (r *NodeRegistry) CreateNodeID(conn net.Conn) uint8 {
	r.mu.Lock()
	defer r.mu.Unlock()

	nodeID := util.GetUniqueID(r.nodes)
	r.nodes[nodeID] = conn
	return nodeID
}

// RemoveNodeID deletes the node ID and its associated connection.
func (r *NodeRegistry) RemoveNodeID(id uint8) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.nodes, id)
}
