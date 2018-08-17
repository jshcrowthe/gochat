package server

import (
	"bufio"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

func handleClient(conn net.Conn, msgs chan<- message) {
	// Prompt client to identify themselves
	writeColorToConn(conn, "Your Name: ")
	reader := bufio.NewReader(conn)

	name, err := readFromReader(reader)

	if err != nil {
		log.Fatal(err)
	}

	// Welcome User
	count := len(getConns())
	writeColorToConn(conn, fmt.Sprintf("Welcome %s! - There are %d other connected users\n", name, count))

	// Add connection to master list of connections and
	// defer connection cleanup
	addConn(conn)
	defer deleteConn(conn)

	// Handle future incoming messages as text
	for {
		text, err := readFromReader(reader)

		if err != nil {
			break
		}

		msgs <- message{
			Author:    name,
			Text:      text,
			Timestamp: time.Now(),
		}
	}
}
