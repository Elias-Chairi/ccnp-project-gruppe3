package util

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	udpservice "github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/server/udpService"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

const searchDuration = 3 * time.Second

// FindServer sends a UDP discovery message to the multicast group and listens for ACK responses.
// It returns a slice of IP addresses of discovered servers.
func FindServer() ([]net.IP, error) {
	// create UDP socket
	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to find an address: %w", err)
	}

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

	serverAddrChan := make(chan net.IP)
	var wg sync.WaitGroup

	go func() {
		time.Sleep(searchDuration)
		_ = conn.Close() // release thread blocked on ReadFromUDP
		wg.Wait()        // wait for processResponse goroutines to finish
		close(serverAddrChan)
	}()

	for {
		buf := make([]byte, 1024)
		n, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break // timeout reached, stop listening
			} else {
				continue // ignore other errors
			}
		}
		wg.Add(1)
		go processResponse(buf[:n], src, serverAddrChan, &wg)
	}

	// wait until channel is closed
	var serverAddrSlice []net.IP
	for val := range serverAddrChan {
		serverAddrSlice = append(serverAddrSlice, val)
	}

	return serverAddrSlice, nil
}

// processResponse processes incoming UDP responses and sends server IPs to the channel.
func processResponse(data []byte, src *net.UDPAddr, serverAddrChan chan<- net.IP, wg *sync.WaitGroup) {
	defer wg.Done()

	// server message is TLV
	t, err := tlv.DecodeTLV(data)
	if err != nil {
		return
	}

	// server message is ACK/ERROR
	msg, err := messages.DecodeAckErrorMessage(t)
	if err != nil {
		return
	}

	// if message is ERROR, ignore
	if msg.IsError() {
		return
	}

	// send server address to channel (it doesn't close until all responses are processed)
	serverAddrChan <- src.IP
}
