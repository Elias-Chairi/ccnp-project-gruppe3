package udpservice

import (
	"log"
	"net"
)

const (
	MULTICAST_ADDR string = "224.0.0.1:9999"
)

func StartUDPService() {
	groupAddr, err := net.ResolveUDPAddr("udp", MULTICAST_ADDR)
	if err != nil {
		log.Fatal("could not resolve udp address")
	}
	conn, err := net.ListenMulticastUDP("udp", nil, groupAddr)
	if err != nil {
		log.Fatal("Error listening on multicast udp")
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		_, senderAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("error reading from udp buffer %w", err)
			continue
		}
		log.Printf("Recieved request from: %v \n", senderAddr)

		_, err = conn.WriteToUDP([]byte("Ack"), senderAddr)
		if err != nil {
			log.Printf("error writing back to sender: %w", err)
		}
	}
}
