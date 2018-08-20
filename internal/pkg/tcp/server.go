package tcp

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

// handleConnections - Continually accepts connections from a
// server and spawns goroutines to handle each connected client
func handleConnections(server net.Listener) {
	// Infinite loop that accepts all new clients
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("Client connected from: %v", conn.RemoteAddr())

		// Handle future interactions with this client
		go handleClient(conn)
	}
}

// Start - Creates the actual chat server
func Start(ip string, port int) {
	// Start the TCP server
	address := fmt.Sprintf("%s:%d", ip, port)
	server, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("TCP Server Listening on %v", address)
	defer server.Close()

	handleConnections(server)
}
