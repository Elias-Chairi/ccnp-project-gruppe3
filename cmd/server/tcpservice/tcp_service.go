package tcpservice

import (
	"fmt"
	"log"
	"net"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
)

var tcpServiceAddress = net.TCPAddr{
	// localhost address
	IP: net.ParseIP("127.0.0.1"),
	// port sat in protocol document
	Port: 6000,
}

func StartTCPService() {

	listner, err := net.ListenTCP("tcp4", &tcpServiceAddress)
	if err != nil {
		log.Fatal("Error listening TCP:", err)
	}
	defer func() {
		_ = listner.Close()
	}()

	log.Printf("Listening on: %s\n", listner.Addr())

	// connChan := make(chan net.Conn)

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go handleRegistration(conn)
	}
}

// Handle connection
func handleRegistration(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading from conn:", err)
		_ = conn.Close()
		return
	}
	err = registerClient(buf[:n], conn)
	if err != nil {
		log.Println("Failed to registe client:", err)
	}
}

func registerClient(msg []byte, conn net.Conn) error {
	t, err := tlv.DecodeTLV(msg)
	if err != nil {
		return fmt.Errorf("Failed to decode TLV: %w", err)
	}
	switch t.Type() {
	case uint8(constants.REGISTER_NODE):
		return handleConn(conn, t.Type(), handleNode)

	// case uint8(constants.REGISTER_CONTROL):
	// 	return handleConn()

	default:
		return fmt.Errorf("Message Registration type '%v' is invalid", t.Type())

	}

}

func handleConn(conn net.Conn, ty uint8, f func(ty uint8, t []tlv.TLV) error) error {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return fmt.Errorf("Error reading conn bytes %w", err)
		}
		ts, err := tlv.DecodeMultipleTLVs(buf[:n])
		if err != nil {
			return fmt.Errorf("Error decoding TLV %w", err)
		}
		err = f(ty, ts)
		if err != nil {
			return fmt.Errorf("Error: %w", err)

		}
	}
}

func handleNode(ty uint8, t []tlv.TLV) error {
	return nil

}
