package server

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/mgutz/ansi"
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

	// Setup file logger
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Unable to open logfile")
	}
	defer file.Close()

	fLog := log.New()
	fLog.Formatter = &log.JSONFormatter{}
	fLog.Out = file

	go handleConnections(server, msgs, newConns, deadConns)

	for {
		select {
		case msg := <-msgs:
			// TODO: Queue up appending message to logfile
			go func(msg message) {
				fLog.WithFields(log.Fields{
					"author":           msg.Author,
					"messageTimestamp": msg.Timestamp,
				}).Info(msg.Text)
			}(msg)

			// Write the message to all active connections
			for conn := range activeConns {
				go func(conn net.Conn, msg message) {
					time := msg.Timestamp.Format("01/02/06 03:04PM")
					prefix := ansi.Color(time+" "+msg.Author+">", "white")
					conn.Write([]byte(prefix + " " + msg.Text + "\n"))
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
