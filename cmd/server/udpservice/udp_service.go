package udpservice

import (
	"log"
	"net"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
)

var MULTICAST_ADDR = net.UDPAddr{
	IP:   net.ParseIP("224.0.0.1"),
	Port: 9999,
}

// StartUDPService starts the UDP service that listens for discovery messages and responds with ACKs.
func StartUDPService() {
	// Create a UDP socket bound to the multicast address
	conn, err := net.ListenMulticastUDP("udp4", nil, &MULTICAST_ADDR)
	if err != nil {
		log.Fatalf("Error listening on multicast udp: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	log.Printf("Listening for multicast messages on %s\n", MULTICAST_ADDR.String())

	// Continuously read from the UDP socket
	for {
		buf := make([]byte, 1024)
		n, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("error reading from udp connection: %v", err)
			continue
		}
		go processRequest(buf[:n], src)
	}
}

// processRequest processes incoming UDP discovery requests and sends ACK responses.
func processRequest(data []byte, src *net.UDPAddr) {
	// client message is TLV
	t, err := encoding.DecodeTLV(data)
	if err != nil {
		log.Printf("error decoding TLV: %v", err)
		return
	}

	// client message is DISCOVERY
	_, err = messages.DecodeDiscoveryMessage(t)
	if err != nil {
		log.Printf("error decoding discovery message: %v", err)
		return
	}

	// create reply connection
	replyConn, err := net.DialUDP("udp4", nil, src)
	if err != nil {
		log.Printf("error dialing back to sender: %v", err)
		return
	}

	// create ack message
	msg := messages.AckSuccessMessage()
	tlv, err := msg.Encode()
	if err != nil {
		log.Printf("error encoding ACK message: %v", err)
		return
	}

	// send ACK message
	_, err = replyConn.Write(tlv.Encode())
	if err != nil {
		log.Printf("error writing ACK message: %v", err)
		return
	}
	_ = replyConn.Close()

	log.Printf("Replied to %v with ACK message\n", src)
}
