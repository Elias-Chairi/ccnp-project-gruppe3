package tcpservice

import (
	"net"
	"slices"
	"sync"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
)

var handleControlPanel messageHandler = func(t *tcpService, msg encoding.TLV) messages.TopLevelMessage {
	switch msg.Type() {
	case uint8(constants.COMMAND): // received command from control panel, forward to node

		id := t.pendingReq.Add()
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

// AddControlPanel adds a control panel connection to the registry.
func (r *ControlPanelRegistry) AddControlPanel(conn net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nodes = append(r.nodes, conn)
}

// RemoveControlPanel removes the specified control panel connection from the registry.
func (r *ControlPanelRegistry) RemoveControlPanel(conn net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	i := slices.Index(r.nodes, conn)
	r.nodes = slices.Delete(r.nodes, i, i+1)
}
