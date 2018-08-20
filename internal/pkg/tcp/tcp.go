package tcp

import (
	"net"

	"github.com/jshcrowthe/gochat/internal/pkg/chat"
	"github.com/mgutz/ansi"
)

// Client - A type for TCP clients
type Client struct {
	conn net.Conn
}

// writeColorToConn - Colored write to a net.Conn
func writeColorToConn(conn net.Conn, s string) {
	txt := ansi.Color(s, "green")
	conn.Write([]byte(txt))
}

// writeToConn - writes a string to a net.Conn
func writeToConn(conn net.Conn, s string) {
	conn.Write([]byte(s))
}

// Write - a method to handle writing to a TCP Client
func (c Client) Write(msg chat.Message) {
	time := msg.Timestamp.Format("01/02/06 03:04PM")
	writeColorToConn(c.conn, time+" "+msg.Nickname+" ("+msg.Email+")> ")
	writeToConn(c.conn, msg.Text+"\n")
}
