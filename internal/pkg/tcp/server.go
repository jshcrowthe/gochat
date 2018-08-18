package tcp

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

// Start - Creates the actual chat server
func Start(ip string, port int, logfile string) {
	// Start the TCP server
	address := fmt.Sprintf("%s:%d", ip, port)
	server, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Server Listening on %v - logs captured at %v", address, logfile)
	defer server.Close()

	handleConnections(server)
}
