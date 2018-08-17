package server

import (
	"bufio"
	"net"
	"strings"

	"github.com/mgutz/ansi"
)

func readFromReader(reader *bufio.Reader) (string, error) {
	text, err := reader.ReadString('\n')
	return strings.TrimSpace(text), err
}

func writeColorToConn(conn net.Conn, s string) {
	txt := ansi.Color(s, "green")
	conn.Write([]byte(txt))
}

func writeToConn(conn net.Conn, s string) {
	conn.Write([]byte(s))
}

// Thread Safe Add
func addConn(conn net.Conn) {
	connsMutex.Lock()
	conns[conn] = true
	connsMutex.Unlock()
}

// Thread Safe Delete
func deleteConn(conn net.Conn) {
	connsMutex.Lock()
	delete(conns, conn)
	connsMutex.Unlock()
}

// Thread Safe Get
func getConns() map[net.Conn]bool {
	connsMutex.Lock()
	defer connsMutex.Unlock()

	return conns
}
