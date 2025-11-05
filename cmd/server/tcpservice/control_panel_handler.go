package tcpservice

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

var handleControlPanel ConnHandler = func(t tlv.TLV) (constants.AckErrorCode, error) {
	switch t.Type() {
	case uint8(constants.COMMAND):
		// todo: handle command
		return constants.ACK_SUCCESS, nil
	default:
		return constants.ERR_INVALID_MESSAGE_TYPE, fmt.Errorf("invalid type %v for control panel handler", t.Type())
	}
}
