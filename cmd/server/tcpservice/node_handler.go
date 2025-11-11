package tcpservice

import (
	"net"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util"
)

// Map to assign node IDs to the net.Conn
var nodeIDsMap = make(map[uint8]net.Conn)

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

// Create and assign a new unique node ID
func CreateNodeID(conn net.Conn) uint8 {
	nodeID := util.GetUniqueID(nodeIDsMap)
	nodeIDsMap[nodeID] = conn
	return nodeID
}

func RemoveNodeID(id uint8) {
	delete(nodeIDsMap, id)
}
