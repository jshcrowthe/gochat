package tcp

import (
	"os"

	"github.com/jshcrowthe/gochat/internal/pkg/chat"
	log "github.com/sirupsen/logrus"
)

func startChat(logfile string) {
	// Setup file logger
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Unable to open logfile")
	}
	defer file.Close()

	fLog := log.New()
	fLog.Formatter = &log.JSONFormatter{}
	fLog.Out = file

	for {
		msg := <-chat.MessagesChan
		// Write message to logfile
		go func(msg chat.Message) {
			fLog.WithFields(log.Fields{
				"email":       msg.Email,
				"nickname":    msg.Nickname,
				"received_at": msg.Timestamp,
			}).Info(msg.Text)
		}(msg)

		// Write the message to all active connections
		for client := range chat.GetClients()[chat.TCP] {
			go func(client chat.Writeable, msg chat.Message) {
				client.Write(msg)
			}(client, msg)
		}
	}
}
