package tcpservice

import (
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
)

var handleNode messageHandler = func(t *encoding.TRLV) messages.AckErrorMessage {
	switch t.TLV.Type() {
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
