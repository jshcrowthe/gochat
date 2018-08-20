# `gochat`

A chat server implementation that supports connections via tcp as well as over websocket

## Setup

_NOTE: There are several CLI commands illustrated below, unless mentioned otherwise, these should be run from the root of this repository._

### Dependencies

This application's dependencies are managed via [`govendor`](https://github.com/kardianos/govendor). Install all dependencies by running:

```shell
$ govendor sync
```

### Installation

After the dependencies have been installed, you can build the application by running:

```shell
$ govendor install +p
```

Alternatively you can build/install the main entrypoint located at `cmd/gochat/main.go`

## Running the Server

If you installed this application via `govendor install` or `go install` illustrated above, you can run the `gochat` command to start the server. 

You can also run the application using `go run` like so:

```shell
$ go run cmd/gochat/main.go
```

## Configuration

There are two ways to configure this server. You can see all CLI flags available by running:

```shell
$ gochat -h
```

All configuration options can be set via CLI flags however, for convenience, the following flags can also be configured via a YAML file:

- `ip`
- `tcp-port`
- `http-port`
- `logfile`

A sample config YAML is available below:

```yaml
ip: "127.0.0.1"
tcp-port: 5000
http-port: 5001
logfile: "logfile.txt"
```

## Sources

All of the following sources were used to some degree to aid in the creation of this application. 

- https://github.com/golang-standards/project-layout
- https://golang.org/doc/effective_go.html
- https://gobyexample.com/
- https://github.com/diltram/go-chat
- https://github.com/dbnegative/go-telnet-chatserver
- https://github.com/scotch-io/go-realtime-chat

In addition to these, all of the packages enumerated in `vendor/vendor.json` (and their documentation) were also consulted

### Third Pary Code

Most of the code was written after study of the above sources and the associated documentation. The web UI (located in `/web`) didn't seem like a critical piece of this exercise and as such is the UI from https://github.com/scotch-io/go-realtime-chat adapted to the messages structure of our application. 