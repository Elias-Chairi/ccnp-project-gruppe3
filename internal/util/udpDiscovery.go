package util

import (
	"fmt"
	"net"
	"time"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/server/udpService"
)

func FindServer() (net.Addr, error) {
	groupAddr, err := net.ResolveUDPAddr("udp", udpservice.MULTICAST_ADDR)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve address: %w", err)
	}
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to find an address: %w", err)
	}
	defer conn.Close()
	_, err = conn.WriteToUDP([]byte("hello new phone who dis?"), groupAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to write to group: %w", err)
	}
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))

	_, serverAddress, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read buffer: %w", err)
	}
	return serverAddress, nil
}
