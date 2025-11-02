package udpservice

import (
	"log"
	"net"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/messages"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

var MULTICAST_ADDR = net.UDPAddr{
	IP:   net.ParseIP("224.0.0.1"),
	Port: 9999,
}

func StartUDPService() {
	// Create a UDP socket bound to the multicast address
	conn, err := net.ListenMulticastUDP("udp4", nil, &MULTICAST_ADDR)
	if err != nil {
		log.Fatal("Error listening on multicast udp")
	}
	defer func (){
		if err := conn.Close(); err != nil{
			log.Println("Failed closing the connection: ", err)
		}
	}()

	log.Printf("Listening for multicast messages on %s\n", MULTICAST_ADDR.String())

	buf := make([]byte, 1024)
	for {
		n, src, err := conn.ReadFromUDP(buf)
		go processRequest(buf, n, src, err)
	}
}

func processRequest(buf []byte, n int, src *net.UDPAddr, err error) {
	// check for read error
	if err != nil {
		log.Printf("error reading from udp connection: %v", err)
		return
	}

	// client message is TLV
	t, err := tlv.DecodeTLV(buf[:n])
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
	msg := messages.NewAckMessage(nil)
	encodedMsg, err := msg.Encode()
	if err != nil {
		log.Printf("error encoding ACK message: %v", err)
		return
	}

	// send ACK message
	_, err = replyConn.Write(encodedMsg)
	if err != nil {
		log.Printf("error writing ACK message: %v", err)
		return
	}
	err = replyConn.Close()
	if err != nil {
		log.Println("Error closing reply conn: ", err)
	}

	log.Printf("Replied to %v with ACK message\n", src)
}
