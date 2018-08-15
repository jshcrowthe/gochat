package app

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
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

	// Captures all of the flags that are supported by this application
	// Some support definition via a yaml file, those are marked with the
	// `altsrc` wrapping
	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "path to YAML `FILE` for config",
		},
		altsrc.NewIntFlag(cli.IntFlag{
			Name:  "port",
			Usage: "`PORT` to listen on",
			Value: 8080,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "ip",
			Usage: "IP address to listen on",
			Value: "127.0.0.1",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "logfile",
			Usage: "output `FILE` for server logs",
			Value: "log.txt",
		}),
	}

	app.Flags = flags

	app.Before = func(c *cli.Context) error {
		if _, err := os.Stat(c.String("config")); os.IsNotExist(err) {
			return nil
		}

		return altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config"))(c)
	}

	app.Action = func(c *cli.Context) error {
		for _, f := range flags {
			log.Println(c.Generic(f.GetName()))
		}
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
