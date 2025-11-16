package tcpservice

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	util "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util/general"
	utilNet "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util/net"
)

var handleNode messageHandler = func(t *tcpService, conn *utilNet.SafeConn, tlv encoding.TLV, reqID *uint16) {
	switch tlv.Type() {
	case uint8(constants.SENSOR_UPDATE):
		msg, err := messages.DecodeSensorUpdateMessage(tlv)
		if err != nil {
			// malformed sensor update message, ignore
			return
		}

		// forward sensor update to all connected control panels
		t.NotifyAllControlPanels(messages.NewSensorUpdateMessage(msg.Sensor))
	case uint8(constants.ACK_ERROR_REQUESTID):
		req, ok := t.pendingReq.Get(*reqID)
		if !ok {
			// unknown request ID, ignore
			return
		}

		switch reqMsg := req.Msg.(type) {
		case *messages.CommandMessage:
			// command response from node to control panel
			switch reqMsg.NodeSelector.Type {
			case constants.SINGLE_NODE:
				// acknowledge to original sender
				_ = writeMessage(req.Sender, &req.ReqID, messages.AckSuccessMessage())

				// forward actuator update to all other control panels
				for _, c := range t.ctrlPanReg.GetAllExcept(req.Sender) {
					_ = writeMessage(c, nil, messages.NewActuatorUpdateMessage(reqMsg.ActuatorSelector, reqMsg.ActuatorState))
				}
			case constants.NODE_LIST:
				// maybe future functionality
			case constants.ALL_NODES:
				// maybe future functionality
			}
		default:
			// unknown original message type, ignore
			return
		}

		t.pendingReq.Remove(*reqID)
	default:
		_ = writeMessage(conn, nil, messages.AckErrorMessage{
			Code: constants.ERR_INVALID_MESSAGE_TYPE,
		})
	}
}

// NodeRegistry manages node IDs and their associated connections.
type NodeRegistry struct {
	mu    sync.RWMutex
	nodes map[uint8]*NodeInfo
}

// NodeInfo holds information about a registered node.
type NodeInfo struct {
	Conn      *utilNet.SafeConn
	Sensors   []entity.Sensor[any]
	Actuators []entity.Actuator[any]
}

// func testNodes() []entity.Node {
// 	nodeA := entity.Node{
// 		ID: 1,
// 		Sensors: []entity.Sensor[any]{
// 			{ID: 1, Type: "Temperature", Unit: "°C", Value: 22.3},
// 			{ID: 2, Type: "Humidity", Unit: "%", Value: 55},
// 			{ID: 3, Type: "CO2", Unit: "ppm", Value: 420},
// 		},
// 		Actuators: []entity.Actuator[any]{
// 			{ID: 1, Type: "Heater", State: "OFF"},
// 			{ID: 2, Type: "Fan", State: "ON"},
// 			{ID: 3, Type: "Sprinkler", State: "OFF"},
// 		},
// 	}

// 	// --- Greenhouse B ---
// 	nodeB := entity.Node{
// 		ID: 2,
// 		Sensors: []entity.Sensor[any]{
// 			{ID: 1, Type: "Temperature", Unit: "°C", Value: 25.7},
// 			{ID: 2, Type: "Humidity", Unit: "%", Value: int32(48)},
// 			{ID: 3, Type: "CO2", Unit: "ppm", Value: int32(390)},
// 		},
// 		Actuators: []entity.Actuator[any]{
// 			{ID: 1, Type: "Heater", Unit: "", State: "OFF"},
// 			{ID: 2, Type: "Fan", Unit: "", State: "OFF"},
// 			{ID: 3, Type: "Sprinkler", Unit: "", State: "ON"},
// 		},
// 	}

// 	// --- Greenhouse C ---
// 	nodeC := entity.Node{
// 		ID: 3,
// 		Sensors: []entity.Sensor[any]{
// 			{ID: 1, Type: "Temperature", Unit: "°C", Value: float32(19.5)},
// 			{ID: 2, Type: "Humidity", Unit: "%", Value: int32(62)},
// 			{ID: 3, Type: "CO2", Unit: "ppm", Value: int32(450)},
// 		},
// 		Actuators: []entity.Actuator[any]{
// 			{ID: 1, Type: "Heater", Unit: "", State: false},
// 			{ID: 2, Type: "Fan", Unit: "", State: true},
// 			{ID: 3, Type: "Sprinkler", Unit: "", State: "OFF"},
// 		},
// 	}

// 	// --- Combine all nodes ---
// 	nodes := []entity.Node{nodeA, nodeB, nodeC}
// 	return nodes
// }

// GetAllNodes returns a slice of all registered nodes.
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
		nodes: make(map[uint8]*NodeInfo),
	}
}

// CreateNodeID assigns a unique node ID and stores the connection.
func (r *NodeRegistry) CreateNodeID(conn *utilNet.SafeConn, sensors []entity.Sensor[any], actuators []entity.Actuator[any]) uint8 {
	r.mu.Lock()
	defer r.mu.Unlock()

	nodeID := util.GetUniqueMapKey(r.nodes)
	r.nodes[nodeID] = &NodeInfo{
		Conn:      conn,
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

var ErrNodeNotFound = errors.New("node not found")

// WriteToNode sends a message to the specified node.
func (r *NodeRegistry) WriteToNode(id uint8, reqID *uint16, msg messages.TopLevelMessage) error {
	r.mu.RLock()
	info, exists := r.nodes[id]
	r.mu.RUnlock()

	if !exists {
		return fmt.Errorf("node %d: %w", id, ErrNodeNotFound)
	}

	return writeMessage(info.Conn, reqID, msg)
}
