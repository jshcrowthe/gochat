package server

import (
	"net"

	log "github.com/sirupsen/logrus"
)

func handleConnections(server net.Listener, msgs chan<- message) {
	// Infinite loop that accepts all new clients
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("Client connected from: %v", conn.RemoteAddr())

		// Handle future interactions with this client
		go handleClient(conn, msgs)
	}
}
