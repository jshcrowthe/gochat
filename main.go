// Stub for gochat server

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

/*
 * Pattern for flags adapted from: https://stackoverflow.com/questions/19761963/flag-command-line-parsing-in-golang
 */
var config string
var port int
var ip string
var logfile string

// ServerConfig - The JSON structure used to store configs for server startup
type ServerConfig struct {
	Port    *int    `json:"port,omitempty"`
	IP      *string `json:"ip,omitempty"`
	LogFile *string `json:"logfile,omitempty"`
}

func handleArgs() {
	if config != "" {
		// Validate path is valid
		// Source: https://gist.github.com/mattes/d13e273314c3b3ade33f
		if _, pathErr := os.Stat(config); os.IsNotExist(pathErr) {
			// If path is invalid print error message and exit with exit code 1
			log.Fatalln("Invalid config path passed to application")
		}

		// If we can't read the config file, quit with an error
		raw, readErr := ioutil.ReadFile(config)
		if readErr != nil {
			log.Fatalln("Could not read config file passed to application")
		}

		var data ServerConfig
		marshalErr := json.Unmarshal(raw, &data)

		if marshalErr != nil {
			log.Fatalln("Config file contained malformed JSON")
		}

		// Preference for arguments: defaults < CLI Args < Config File Values
		// Override the existing values if the config exists
		if data.Port != nil {
			port = *data.Port
		}

		if data.IP != nil {
			ip = *data.IP
		}

		if data.LogFile != nil {
			logfile = *data.LogFile
		}
	}
}

func drainLogQueue() {
	fmt.Println("Writing queued logs to disk")
}

func init() {
	// Register all of the supported command line arguments
	flag.StringVar(&config, "config", "", "Path to a JSON file used to configure the server")
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.StringVar(&ip, "ip", "127.0.0.1", "IP address to listen on")
	flag.StringVar(&logfile, "logfile", "log.txt", "Location of logfile")
}

func main() {
	// Parse/Handle all command line args
	flag.Parse()
	handleArgs()

	// Log config
	fmt.Printf("config: %v\r\nport: %v\r\nip: %v\r\nlogfile: %v\r\n", config, port, ip, logfile)

	// Setup ticker to drain log queue and write it to disk
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			drainLogQueue()
		}
	}()

	// Start connection listener
	listener, listenErr := net.Listen("tcp", ip+":"+strconv.Itoa(port))
	if listenErr != nil {
		log.Fatalf("Could not setup listener at %v:%v\r\nError: %v", ip, port, listenErr)
	}

	log.Println("listening on: ", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("could not accept connection %v ", err)
		}
		//create new client on connection
		fmt.Println(conn)
	}
}
