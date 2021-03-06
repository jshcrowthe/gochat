package app

import (
	"os"

	"github.com/jshcrowthe/gochat/internal/pkg/chat"
	"github.com/jshcrowthe/gochat/internal/pkg/tcp"
	"github.com/jshcrowthe/gochat/internal/pkg/websocket"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

const (
	name    = "gochat"
	usage   = "Telnet chat server written in golang"
	version = "1.0.0"
)

var app *cli.App

func init() {
	// Creates a new cli application
	app = cli.NewApp()

	// Define information about the application
	app.Version = version
	app.Name = name
	app.Usage = usage

	// Captures all of the flags that are supported by this application
	// Some support definition via a yaml file, those are marked with the
	// `altsrc` wrapping
	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "path to YAML `FILE` for config",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "set debug mode for the application",
		},
		altsrc.NewIntFlag(cli.IntFlag{
			Name:  "tcp-port",
			Usage: "`PORT` to start TCP listener on",
			Value: 8080,
		}),
		altsrc.NewIntFlag(cli.IntFlag{
			Name:  "http-port",
			Usage: "`PORT` to start HTTP listener on (also used for WebSocket implementation)",
			Value: 8081,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "ip",
			Usage: "IP address to listen on",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "logfile",
			Usage: "output `FILE` for server logs",
			Value: "log.txt",
		}),
	}

	app.Flags = flags

	app.Before = func(c *cli.Context) error {
		// If the --debug flag is passed set the log level appropriately
		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		if _, err := os.Stat(c.String("config")); os.IsNotExist(err) {
			log.Debug("No config file passed")
			return nil
		}

		return altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config"))(c)
	}

	app.Action = func(c *cli.Context) error {
		// Starts TCP connection handler
		go tcp.Start(c.String("ip"), c.Int("tcp-port"))

		// Starts websocket connection handler
		go websocket.Start(c.String("ip"), c.Int("http-port"))

		// Start chat
		chat.Start(c.String("logfile"))
		return nil
	}
}

// Run - Starts the server app and handles all of the flag parsing
func Run() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
