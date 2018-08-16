package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func readFromReader(reader *bufio.Reader) (string, error) {
	text, err := reader.ReadString('\n')
	return strings.TrimSpace(text), err
}

func handleClient(conn net.Conn, messages chan<- message, deadConns chan<- net.Conn) {
	// Prompt client to identify themselves
	conn.Write([]byte("Your Name: "))
	reader := bufio.NewReader(conn)

	author, err := readFromReader(reader)

	if err != nil {
		log.Fatal(err)
	}

	// Welcome User
	conn.Write([]byte(fmt.Sprintf("Welcome %s!\n", author)))

	// Handle future incoming messages as text
	for {
		text, err := readFromReader(reader)

		if err != nil {
			break
		}

		messages <- message{
			Author:    author,
			Text:      text,
			Timestamp: time.Now(),
		}
	}

	// Once loop has broken, consider connection closed
	deadConns <- conn
}

func handleConnections(server net.Listener, messages chan<- message, newConns chan<- net.Conn, deadConns chan<- net.Conn) {
	// Infinite loop that accepts all new clients
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("Client connected from: %v", conn.RemoteAddr())

		// Announce the new connection to newConns channel
		newConns <- conn

		// Handle future interactions with this client
		go handleClient(conn, messages, deadConns)
	}
}
