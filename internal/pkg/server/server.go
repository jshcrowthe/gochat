package server

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/mgutz/ansi"
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

var msgs = make(chan message)

// Start - Creates the actual chat server
func Start(ip string, port int, logfile string) {
	// Start the TCP server
	server, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Server Listening on %v:%v - logs captured at %v", ip, port, logfile)
	defer server.Close()

	// Setup file logger
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Unable to open logfile")
	}
	defer file.Close()

	fLog := log.New()
	fLog.Formatter = &log.JSONFormatter{}
	fLog.Out = file

	go handleConnections(server)

	for {
		msg := <-msgs
		// TODO: Queue up appending message to logfile
		go func(msg message) {
			fLog.WithFields(log.Fields{
				"author":           msg.Author,
				"messageTimestamp": msg.Timestamp,
			}).Info(msg.Text)
		}(msg)

		// Write the message to all active connections
		for conn := range conns {
			go func(conn net.Conn, msg message) {
				time := msg.Timestamp.Format("01/02/06 03:04PM")
				prefix := ansi.Color(time+" "+msg.Author+">", "white")
				conn.Write([]byte(prefix + " " + msg.Text + "\n"))
			}(conn, msg)
		}
	}
}
