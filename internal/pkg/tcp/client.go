package tcp

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/jshcrowthe/gochat/internal/pkg/chat"
	log "github.com/sirupsen/logrus"
)

func readFromReader(reader *bufio.Reader) (string, error) {
	text, err := reader.ReadString('\n')
	return strings.TrimSpace(text), err
}

func handleClient(conn net.Conn) {
	// Prompt client to identify themselves
	writeColorToConn(conn, "Your Name: ")
	reader := bufio.NewReader(conn)

	name, err := readFromReader(reader)

	if err != nil {
		log.Fatal(err)
	}

	// Welcome User
	count := len(chat.GetClients()[chat.TCP])
	writeColorToConn(conn, fmt.Sprintf("Welcome %s! - There are %d other connected users\n", name, count))

	// Add connection to master list of connections and
	// defer connection cleanup
	client := &Client{
		conn: conn,
	}
	chat.AddClient(chat.TCP, client)
	defer chat.DeleteClient(chat.TCP, client)

	// Handle future incoming messages as text
	for {
		text, err := readFromReader(reader)

		if err != nil {
			break
		}

		chat.MessagesChan <- chat.Message{
			Nickname:  name,
			Text:      text,
			Timestamp: time.Now(),
		}
	}
}
