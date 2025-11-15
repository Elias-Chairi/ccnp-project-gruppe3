package udpservice

import (
	"log"
	"net"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/encoding"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util"
)

// StartUDPService starts the UDP service that listens for discovery messages and responds with ACKs.
func StartUDPService() {
	// Create a UDP socket bound to the multicast address
	conn, err := net.ListenMulticastUDP("udp4", nil, &util.MULTICAST_ADDR)
	if err != nil {
		log.Fatalf("Error listening on multicast udp: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	log.Printf("Listening for multicast messages on %s\n", util.MULTICAST_ADDR.String())

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
	defer replyConn.Close() //Ensures the UDP socket always closes, even if Write or Decode fails.

	// create ack message
	msg := messages.AckSuccessMessage()
	tlv, err := msg.Encode()
	if err != nil {
		log.Printf("error encoding ACK message: %v", err)
		return
	}

	// Send ACK message
	_, err = replyConn.Write(tlv.Encode())
	if err != nil {
		log.Printf("error writing ACK message: %v", err)
		return
	}

	log.Printf("Replied to %v with ACK message\n", src)
}
