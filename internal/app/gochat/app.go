package app

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
)

const (
	appName    = "gochat"
	appUsage   = "Telnet chat server written in golang"
	appVersion = "1.0.0"
)

var app *cli.App

func init() {
	// Creates a new cli application
	app = cli.NewApp()

	// Define information about the application
	app.Version = appVersion
	app.Name = appName
	app.Usage = appUsage

	// Define a fxn to use if cmd isn't recognized
	app.CommandNotFound = func(ctx *cli.Context, cmd string) {
		log.Printf("unknown command - %v \n\n", cmd)
		cli.ShowAppHelp(ctx)
	}

	app.Flags = []cli.Flag{
		altsrc.NewIntFlag(cli.IntFlag{
			Name:  "port, p",
			Usage: "port to listen on",
			Value: 8080,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "ip, i",
			Usage: "IP address to listen on",
			Value: "127.0.0.1",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "logfile, l",
			Usage: "output path for server logs",
			Value: "log.txt",
		}),
	}

	app.Before = altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config, c"))
}

// Run - Starts the server app and handles all of the flag parsing
func Run() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
