package tcpservice

import (
	"errors"
	"fmt"
	"slices"
	"sync"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	utilNet "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util/net"
)

var handleControlPanel messageHandler = func(t *tcpService, conn *utilNet.SafeConn, tlv encoding.TLV, senderReqID *uint16) {
	switch tlv.Type() {
	case uint8(constants.COMMAND): // received command from control panel, forward to node
		fmt.Println("RECEIVED COMMAND MESSAGE FROM CONTROL PANEL")
		msg, err := messages.DecodeCommandMessage(tlv, true)
		if err != nil {
			_ = writeMessage(conn, nil, messages.AckErrorRequestIDMessage{
				Code: constants.ERR_MALFORMED_MESSAGE,
			})
			return
		}

		reqID := t.pendingReq.Add(utilNet.Request{
			Msg:    msg,
			Sender: conn,
			ReqID:  *senderReqID,
		})

		switch msg.NodeSelector.Type {
		case constants.SINGLE_NODE:
			// because the node don't need to know about the node selector, we create a copy without it
			msgToNode := new(messages.CommandMessage)
			*msgToNode = *msg
			msgToNode.NodeSelector = nil

			fmt.Println("FORWARDING COMMAND MESSAGE TO NODE", msg.NodeSelector.NodeIDs[0])
			err = t.nodeReg.WriteToNode(msg.NodeSelector.NodeIDs[0], &reqID, msgToNode)
			if err != nil {
				if errors.Is(err, ErrNodeNotFound) {
					_ = writeMessage(conn, nil, messages.AckErrorRequestIDMessage{
						Code: constants.ERR_UNKNOWN_NODE_ID,
					})
				} else {
					_ = writeMessage(conn, nil, messages.AckErrorRequestIDMessage{
						Code: constants.ERR_INTERNAL_SERVER_ERROR,
					})
				}
				t.pendingReq.Remove(reqID)
				return
			}
		case constants.NODE_LIST:
			// maybe future functionality
		case constants.ALL_NODES:
			// maybe future functionality
		}

	default:
		_ = writeMessage(conn, nil, messages.AckErrorRequestIDMessage{
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
