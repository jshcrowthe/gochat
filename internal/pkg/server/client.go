package server

import (
	"bufio"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

func handleClient(conn net.Conn) {
	// Prompt client to identify themselves
	conn.Write([]byte("Your Name: "))
	reader := bufio.NewReader(conn)

	author, err := readFromReader(reader)

	if err != nil {
		log.Fatal(err)
	}

	// Welcome User
	conn.Write([]byte(fmt.Sprintf("Welcome %s!\n", author)))

	connsMutex.Lock()
	conns[conn] = true
	connsMutex.Unlock()

	// Handle future incoming messages as text
	for {
		text, err := readFromReader(reader)

		if err != nil {
			break
		}

		msgs <- message{
			Author:    author,
			Text:      text,
			Timestamp: time.Now(),
		}
	}

	// Once loop breaks, connection is closed
	connsMutex.Lock()
	delete(conns, conn)
	connsMutex.Unlock()
}
