package server

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type message struct {
	Author    string
	Text      string
	Timestamp time.Time
}

// Start - Creates the actual chat server
func Start(ip string, port int, logfile string) {
	// Start the TCP server
	server, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Server Listening on %v:%v - logs captured at %v", ip, port, logfile)

	// Channels for communication across goroutines
	msgs := make(chan message)
	newConns := make(chan net.Conn)
	deadConns := make(chan net.Conn)

	// Slice to hold all active connections
	connCount := 0
	activeConns := make(map[net.Conn]int)

	go handleConnections(server, msgs, newConns, deadConns)

	for {
		select {
		case msg := <-msgs:
			// TODO: Queue up appending message to logfile

			for conn := range activeConns {
				go func(conn net.Conn, msg message) {
					conn.Write([]byte(`\33[2K`))
					conn.Write([]byte(msg.Author + ": " + msg.Text + "\n"))
				}(conn, msg)
			}
		case conn := <-newConns:
			activeConns[conn] = connCount
			connCount++
		case conn := <-deadConns:
			delete(activeConns, conn)
		}
	}
}
