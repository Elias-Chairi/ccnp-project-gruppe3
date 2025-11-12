package tcpservice

import (
	"net"
	"sync"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
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

func GetUniqueID(existingIDs map[uint8]NodeInfo) uint8 {
	newID := uint8(len(existingIDs) + 1)
	for {
		if _, exists := existingIDs[newID]; !exists {
			break
		}
		newID++
	}
	return newID
}

// NodeRegistry manages node IDs and their associated connections.
type NodeRegistry struct {
	mu    sync.RWMutex
	nodes map[uint8]NodeInfo
}

// NodeInfo holds information about a registered node.
type NodeInfo struct {
	conn      net.Conn
	Sensors   []entity.Sensor[any]
	Actuators []entity.Actuator[any]
}

func (r *NodeRegistry) GetAllNodes() []entity.Node {
	r.mu.Lock()
	defer r.mu.Unlock()

	nodes := make([]entity.Node, len(r.nodes))
	i := 0
	for id, info := range r.nodes {
		// create entity.Node from NodeInfo
		// registry holds data about each node
		nodes[i] = entity.Node{
			ID:        id,
			Sensors:   info.Sensors,
			Actuators: info.Actuators,
		}
		i++
	}
	return nodes
}

// NewNodeRegistry initializes and returns a new NodeRegistry.
func NewNodeRegistry() *NodeRegistry {
	return &NodeRegistry{
		nodes: make(map[uint8]NodeInfo),
	}
}

// CreateNodeID assigns a unique node ID and stores the connection.
func (r *NodeRegistry) CreateNodeID(conn net.Conn, sensors []entity.Sensor[any], actuators []entity.Actuator[any]) uint8 {
	r.mu.Lock()
	defer r.mu.Unlock()

	nodeID := GetUniqueID(r.nodes)
	r.nodes[nodeID] = NodeInfo{
		conn:      conn,
		Sensors:   sensors,
		Actuators: actuators,
	}
	return nodeID
}

// RemoveNodeID deletes the node ID and its associated connection.
func (r *NodeRegistry) RemoveNodeID(id uint8) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.nodes, id)
}
