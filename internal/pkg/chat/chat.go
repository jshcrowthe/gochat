package chat

import (
	"sync"
	"time"
)

// Message - A struct that is used to describe all messages passed through the server
type Message struct {
	Email     string
	Text      string
	Nickname  string
	Timestamp time.Time
}

// MessagesChan - A channel that is the primary message stream for the application
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

// Start - Starts the chat handling process
func Start() {
	for {
		msg := <-MessagesChan

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
