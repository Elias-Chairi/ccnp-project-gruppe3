package util

import (
	"context"
	"fmt"
	"net"
	"time"

	udpservice "github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/server/udpService"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

const searchDuration = 3 * time.Second

func FindServer() ([]net.IP, error) {

	// create UDP socket
	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to find an address: %w", err)
	}

	defer func (){
		_ = conn.Close();
	}()
	// send discovery message to multicast group
	msg := messages.NewDiscoveryMessage()
	msgEncoded, err := msg.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode discovery message: %w", err)
	}
	_, err = conn.WriteToUDP(msgEncoded, &udpservice.MULTICAST_ADDR)
	if err != nil {
		return nil, fmt.Errorf("failed to write to group: %w", err)
	}

	var serverAddresses []net.IP
	buf := make([]byte, 1024)

	// context with search duration
	ctx, cancel := context.WithTimeout(context.Background(), searchDuration)
	defer cancel()
	go func() {
		for {
			n, src, err := conn.ReadFromUDP(buf)
			go processResponse(buf, n, src, err, &serverAddresses)
		}
	}()
	<-ctx.Done()

	return serverAddresses, nil
}

func processResponse(buf []byte, n int, src *net.UDPAddr, err error, serverAddresses *[]net.IP) {
	// check for read error
	if err != nil {
		return
	}

	// client message is TLV
	t, err := tlv.DecodeTLV(buf[:n])
	if err != nil {
		return
	}

	// client message is ACK/ERROR
	msg, err := messages.DecodeAckErrorMessage(t)
	if err != nil {
		return
	}

	// if message is ERROR, ignore
	if msg.IsError() {
		return
	}

	*serverAddresses = append(*serverAddresses, src.IP)
}
