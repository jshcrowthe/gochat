package websocket

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
	"github.com/jshcrowthe/gochat/internal/pkg/chat"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	client := &Client{
		conn: ws,
	}

	chat.AddClient(chat.WS, client)
	defer chat.DeleteClient(chat.WS, client)

	for {
		var msg chat.Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			break
		}
		// Send the newly received message to the broadcast channel
		chat.MessagesChan <- msg
	}

}

// Start - Registers HTTP handler
func Start(ip string, port int) {
	// filesystem handler for the web app
	fs := http.FileServer(http.Dir("../../public"))
	http.Handle("/", fs)

	// Websocket endpoint
	http.HandleFunc("/ws", handleConnections)

	address := fmt.Sprintf("%v:%v", ip, port)

	log.Infof("HTTP Server Listening on %v", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
