package tcp

import (
	"fmt"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type message struct {
	Author    string
	Text      string
	Timestamp time.Time
}

var (
	conns      = make(map[net.Conn]bool)
	connsMutex sync.Mutex
)

// Start - Creates the actual chat server
func Start(ip string, port int, logfile string) {
	// Start the TCP server
	address := fmt.Sprintf("%s:%d", ip, port)
	server, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Server Listening on %v - logs captured at %v", address, logfile)
	defer server.Close()

	msgs := make(chan message)

	go handleConnections(server, msgs)

	// This is the main "keep-alive" process
	startChat(msgs, logfile)
}
