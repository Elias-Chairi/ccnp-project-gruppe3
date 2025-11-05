package tcpservice

import (
	"fmt"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

var handleNode ConnHandler = func(t tlv.TLV) (constants.AckErrorCode, error) {
	switch t.Type() {
	case uint8(constants.SENSOR_UPDATE):
		// todo: handle sensor update
		return constants.ACK_SUCCESS, nil
	case uint8(constants.ACK_ERROR):
		// todo: handle ack error
		return constants.ACK_SUCCESS, nil
	default:
		return constants.ERR_INVALID_MESSAGE_TYPE, fmt.Errorf("invalid type %v for node handler", t.Type())
	}
}
