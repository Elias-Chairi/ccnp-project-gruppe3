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
	defer func() {
		_ = conn.Close()
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

	// create timeout context to stop searching after searchDuration
	ctx, cancel := context.WithTimeout(context.Background(), searchDuration)
	var serverAddresses []net.IP
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				buf := make([]byte, 1024)
				n, src, err := conn.ReadFromUDP(buf)
				if err != nil {
					continue
				}
				go processResponse(buf[:n], src, &serverAddresses)
			}
		}
	}()
	<-ctx.Done()
	_ = conn.Close() // release thread blocked on ReadFromUDP

	return serverAddresses, nil
}

func processResponse(data []byte, src *net.UDPAddr, serverAddresses *[]net.IP) {
	// client message is TLV
	t, err := tlv.DecodeTLV(data)
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
