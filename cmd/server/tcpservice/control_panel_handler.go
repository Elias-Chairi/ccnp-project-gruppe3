package tcpservice

import (
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
)

var handleControlPanel messageHandler = func(t *encoding.TRLV) messages.AckErrorMessage {
	switch t.TLV.Type() {
	case uint8(constants.COMMAND):
		// todo: handle command
		return messages.AckSuccessMessage()
	default:
		return messages.AckErrorMessage{
			Code: constants.ERR_INVALID_MESSAGE_TYPE,
		}
	}
}
