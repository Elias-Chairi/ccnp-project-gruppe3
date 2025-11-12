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

type NodeInfo struct {
	conn      net.Conn
	Sensors   []entity.Sensor[any]
	Actuators []entity.Actuator[any]
}

func (r *NodeRegistry) GetAllNodes() []entity.Node {
	r.mu.Lock()
	defer r.mu.Unlock()

	nodes := make([]entity.Node, len(r.nodes))
	for i := range r.nodes {
		nodes[i] = entity.Node{
			ID:        i,
			Sensors:   listOfValuesToListOfPointers(r.nodes[i].Sensors),
			Actuators: listOfValuesToListOfPointers(r.nodes[i].Actuators),
		}
	}
	return nodes
}

func listOfValuesToListOfPointers[T any](values []T) []*T {
	pointers := make([]*T, len(values))
	for i := range values {
		pointers[i] = &values[i]
	}
	return pointers
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
