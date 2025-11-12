package tcpservice

import (
	"net"
	"slices"
	"sync"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
)

var handleControlPanel messageHandler = func(msg *encoding.Message) messages.AckErrorMessage {
	switch msg.TLV.Type() {
	case uint8(constants.COMMAND):
		// todo: handle command
		return messages.AckSuccessMessage()
	default:
		return messages.AckErrorMessage{
			Code: constants.ERR_INVALID_MESSAGE_TYPE,
		}
	}
}

// NodeRegistry manages node IDs and their associated connections.
type ControlPanelRegistry struct {
	mu    sync.RWMutex
	nodes []net.Conn
}

// CreateNodeID assigns a unique node ID and stores the connection.
func (r *ControlPanelRegistry) AddControlPanel(conn net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nodes = append(r.nodes, conn)
}

// RemoveNodeID deletes the node ID and its associated connection.
func (r *ControlPanelRegistry) RemoveControlPanel(conn net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	i := slices.Index(r.nodes, conn)
	r.nodes = slices.Delete(r.nodes, i, i+1)
}
