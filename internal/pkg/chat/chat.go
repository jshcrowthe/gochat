package chat

import (
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Message - A struct that is used to describe all messages passed through the server
type Message struct {
	Email     string    `json:"email"`
	Text      string    `json:"text"`
	Nickname  string    `json:"nickname"`
	Timestamp time.Time `json:"timestamp"`
}

// MessagesChan - T primary message stream for the application
// NOTE: To support more than one channel convert you could
// refactor this from a `chan Message` to a
// `map[string]chan Message` and track the active channel
// for each connected client
var MessagesChan = make(chan Message)

// Writeable - A writeable connection
type Writeable interface {
	Write(msg Message)
}

// ConnectionType - A type for all persistent connection types
type ConnectionType int

// All supported ConnectionType enum values
const (
	TCP ConnectionType = iota
	WS
)

// Clients - All persistent connections
var (
	mutex   sync.Mutex
	clients = make(map[ConnectionType]map[Writeable]bool)
)

// AddClient - Adds a client to the client map
func AddClient(t ConnectionType, w Writeable) {
	mutex.Lock()

	// Initialize map if it hasn't already
	if clients[t] == nil {
		clients[t] = make(map[Writeable]bool)
	}

	clients[t][w] = true
	mutex.Unlock()
}

// DeleteClient - Deletes a client from the client map
func DeleteClient(t ConnectionType, w Writeable) {
	mutex.Lock()
	delete(clients[t], w)
	mutex.Unlock()
}

// GetClients - Gets the master client list
func GetClients() map[ConnectionType]map[Writeable]bool {
	mutex.Lock()
	defer mutex.Unlock()

	return clients
}

// GetClientCount - Gets the count of all connected clients
func GetClientCount() int {
	mutex.Lock()
	defer mutex.Unlock()

	sum := 0

	for _, c := range clients {
		if clients == nil {
			continue
		}

		for clients := range c {
			_ = clients
			sum++
		}
	}

	return sum
}

// Start - Starts the chat handling process
func Start(logfile string) {
	// Setup a file logger
	var m sync.Mutex
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Unable to open logfile")
	}
	defer file.Close()

	fileLog := log.New()
	fileLog.Formatter = &log.JSONFormatter{}
	fileLog.Out = file

	// Infinite loop to keep process alive
	for {

		// Receive a message from the `MessagesChan`
		msg := <-MessagesChan

		// Write the message as a log to disk
		go func(msg Message) {
			m.Lock()
			fileLog.WithFields(log.Fields{
				"email":       msg.Email,
				"nickname":    msg.Nickname,
				"received_at": msg.Timestamp,
			}).Info(msg.Text)
			m.Unlock()
		}(msg)

		// Loop through all connected clients, of all types, and call
		// their respective `write` functions
		for _, c := range clients {
			if clients == nil {
				continue
			}

			for client := range c {
				go client.Write(msg)
			}
		}
	}
}
