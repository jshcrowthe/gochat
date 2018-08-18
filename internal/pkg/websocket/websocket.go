package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/jshcrowthe/gochat/internal/pkg/chat"
)

// Client - A type for TCP clients
type Client struct {
	conn *websocket.Conn
}

// Write - A write implementation for websockets
func (c Client) Write(msg chat.Message) {
	c.conn.WriteJSON(msg)
}
