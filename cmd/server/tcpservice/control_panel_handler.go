package tcpservice

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

func handleControlPanel(t tlv.TLV) error {
	switch t.Type() {
	case uint8(constants.COMMAND):
		// todo: handle command
		return nil
	case uint8(constants.ACK_ERROR):
		// todo: handle ack error
		return nil
	default:
		return fmt.Errorf("invalid type %v for control panel handler", t.Type())
	}
}
