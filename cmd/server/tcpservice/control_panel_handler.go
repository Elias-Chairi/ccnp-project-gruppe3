package tcpservice

import (
	"errors"
	"slices"
	"sync"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	utilNet "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util/net"
)

var handleControlPanel messageHandler = func(t *tcpService, conn *utilNet.SafeConn, tlv encoding.TLV, reqID *uint16) {
	switch tlv.Type() {
	case uint8(constants.COMMAND): // received command from control panel, forward to node
		msg, err := messages.DecodeCommandMessage(tlv, true)
		if err != nil {
			_ = writeMessage(conn, nil, messages.AckErrorMessage{
				Code: constants.ERR_MALFORMED_MESSAGE,
			})
			return
		}

		reqID := t.pendingReq.Add(Request{
			Msg:    msg,
			Sender: conn,
		})

		switch msg.NodeSelector.Type {
		case constants.SINGLE_NODE:
			err = t.nodeReg.WriteToNode(msg.NodeSelector.NodeIDs[0], &reqID, msg)
			if err != nil {
				if errors.Is(err, ErrNodeNotFound) {
					_ = writeMessage(conn, nil, messages.AckErrorMessage{
						Code: constants.ERR_UNKNOWN_NODE_ID,
					})
				} else {
					_ = writeMessage(conn, nil, messages.AckErrorMessage{
						Code: constants.ERR_INTERNAL_SERVER_ERROR,
					})
				}
				return
			}
		case constants.NODE_LIST:
			// maybe future functionality
		case constants.ALL_NODES:
			// maybe future functionality
		}

	default:
		writeMessage(conn, nil, messages.AckErrorMessage{
			Code: constants.ERR_INVALID_MESSAGE_TYPE,
		})
	}
}

// NodeRegistry manages node IDs and their associated connections.
type ControlPanelRegistry struct {
	mu      sync.RWMutex
	ctrlPan []*utilNet.SafeConn
}

// GetAll returns all control panel connections.
func (r *ControlPanelRegistry) GetAll() []*utilNet.SafeConn {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.ctrlPan
}

// GetAllExcept returns all control panel connections except the specified one.
func (r *ControlPanelRegistry) GetAllExcept(except *utilNet.SafeConn) []*utilNet.SafeConn {
	var conns []*utilNet.SafeConn
	for _, c := range r.GetAll() {
		if c != except {
			conns = append(conns, c)
		}
	}
	return conns
}

// AddControlPanel adds a control panel connection to the registry.
func (r *ControlPanelRegistry) AddControlPanel(conn *utilNet.SafeConn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.ctrlPan = append(r.ctrlPan, conn)
}

// RemoveControlPanel removes the specified control panel connection from the registry.
func (r *ControlPanelRegistry) RemoveControlPanel(conn *utilNet.SafeConn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	i := slices.Index(r.ctrlPan, conn)
	r.ctrlPan = slices.Delete(r.ctrlPan, i, i+1)
}
