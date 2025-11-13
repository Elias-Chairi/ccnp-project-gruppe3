package connectionhandler

import (
	"fmt"
	"net"
	"time"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
)

type ConnectionHandler struct {
	conn      net.Conn
	ServerIPs []net.IP
}

func (c *ConnectionHandler) Connect() error {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: c.ServerIPs[0], Port: 6000})
	if err != nil {
		return fmt.Errorf("Failed to connect: %w", err)
	}
	c.conn = conn
	return nil
}

func (c *ConnectionHandler) Register() (*[]entity.Node, error) {
	msg := messages.NewRegisterControlMessage()
	encoded, err := msg.Encode()
	if err != nil {
		return nil, fmt.Errorf("Failed to encode: %w", err)
	}

	_, err = c.conn.Write(encoded)
	if err != nil {
		return nil, fmt.Errorf("Failed to write to conn %w", err)
	}

	readmsg, err := encoding.ReadNextMessage(c.conn, time.Second*10)

	if err != nil {
		return nil, fmt.Errorf("Failed to read next message: %w", err)
	}

	decodedmsg, err := messages.DecodeAckErrorMessage(readmsg.TLV)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode ack error message: %w", err)
	}

	if decodedmsg.IsError() {
		return nil, fmt.Errorf("Errorcode %d: register response is error: %s", decodedmsg.Code, decodedmsg.Data)
	}

	tlvs, err := encoding.DecodeMultipleTLVs([]byte(decodedmsg.Data))
	if err != nil {
		return nil, fmt.Errorf("Failed to decode message: %w", err)
	}

	var nodes []entity.Node

	for _, tlv := range tlvs {
		node, err := encoding.DecodeNodeEntry(tlv)
		if err != nil {
			return nil, fmt.Errorf("Failed to decode node entry: %w", err)
		}
		nodes = append(nodes, *node)
	}

	return &nodes, nil
}
