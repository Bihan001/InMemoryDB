package main

import (
	"flag"

	"github.com/Bihan001/MyDB/config"
	"github.com/Bihan001/MyDB/server"
)

func initFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "Server Host")
	flag.IntVar(&config.Port, "port", 6389, "Server Port")
	flag.IntVar(&config.MaxClients, "maxClients", 20000, "Max Clients to support")
	flag.Parse()
}

func main() {
	initFlags()
	var server server.Server = &server.AsyncTCPServer{}
	server.RunServer()
}
