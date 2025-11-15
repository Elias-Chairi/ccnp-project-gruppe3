package tcpservice

import (
	"net"
	"sync"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util"
)

var handleNode messageHandler = func(t *tcpService, msg encoding.TLV) messages.TopLevelMessage {
	switch msg.Type() {
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
	nodes map[uint8]NodeInfo
}

// NodeInfo holds information about a registered node.
type NodeInfo struct {
	conn      net.Conn
	Sensors   []entity.Sensor[any]
	Actuators []entity.Actuator[any]
}

func testNodes() []entity.Node {
	nodeA := entity.Node{
		ID: 1,
		Sensors: []entity.Sensor[any]{
			{ID: 1, Type: "Temperature", Unit: "°C", Value: 22.3},
			{ID: 2, Type: "Humidity", Unit: "%", Value: 55},
			{ID: 3, Type: "CO2", Unit: "ppm", Value: 420},
		},
		Actuators: []entity.Actuator[any]{
			{ID: 1, Type: "Heater", State: "OFF"},
			{ID: 2, Type: "Fan", State: "ON"},
			{ID: 3, Type: "Sprinkler", State: "OFF"},
		},
	}

	// --- Greenhouse B ---
	nodeB := entity.Node{
		ID: 2,
		Sensors: []entity.Sensor[any]{
			{ID: 1, Type: "Temperature", Unit: "°C", Value: 25.7},
			{ID: 2, Type: "Humidity", Unit: "%", Value: int32(48)},
			{ID: 3, Type: "CO2", Unit: "ppm", Value: int32(390)},
		},
		Actuators: []entity.Actuator[any]{
			{ID: 1, Type: "Heater", Unit: "", State: "OFF"},
			{ID: 2, Type: "Fan", Unit: "", State: "OFF"},
			{ID: 3, Type: "Sprinkler", Unit: "", State: "ON"},
		},
	}

	// --- Greenhouse C ---
	nodeC := entity.Node{
		ID: 3,
		Sensors: []entity.Sensor[any]{
			{ID: 1, Type: "Temperature", Unit: "°C", Value: float32(19.5)},
			{ID: 2, Type: "Humidity", Unit: "%", Value: int32(62)},
			{ID: 3, Type: "CO2", Unit: "ppm", Value: int32(450)},
		},
		Actuators: []entity.Actuator[any]{
			{ID: 1, Type: "Heater", Unit: "", State: false},
			{ID: 2, Type: "Fan", Unit: "", State: true},
			{ID: 3, Type: "Sprinkler", Unit: "", State: "OFF"},
		},
	}

	// --- Combine all nodes ---
	nodes := []entity.Node{nodeA, nodeB, nodeC}
	return nodes
}

// GetAllNodes returns a slice of all registered nodes.
func (r *NodeRegistry) GetAllNodes() []entity.Node {
	r.mu.Lock()
	defer r.mu.Unlock()

	return testNodes()

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

	nodeID := util.GetUniqueMapKey(r.nodes)
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
