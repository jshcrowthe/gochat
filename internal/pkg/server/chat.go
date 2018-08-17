package server

import (
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

func startChat(msgs <-chan message, logfile string) {
	// Setup file logger
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Unable to open logfile")
	}
	defer file.Close()

	fLog := log.New()
	fLog.Formatter = &log.JSONFormatter{}
	fLog.Out = file

	for {
		msg := <-msgs
		// Write message to logfile
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
				writeColorToConn(conn, time+" "+msg.Author+"> ")
				writeToConn(conn, msg.Text+"\n")
			}(conn, msg)
		}
	}
}