package main

import (
    "flag"

    "github.com/Bihan001/MyDB/internal/config"
    "github.com/Bihan001/MyDB/internal/engine"
    "github.com/Bihan001/MyDB/internal/server"
)

func setupFlags() {
    flag.StringVar(&config.ServerHost, "host", "0.0.0.0", "Server Host")
    flag.IntVar(&config.ServerPort, "port", 6389, "Server Port")
    flag.IntVar(&config.ConnectionLimit, "maxClients", 20000, "Max Clients to support")
    flag.Parse()
}

func main() {
    setupFlags()
    var runner server.ServiceRunner = server.NewAsyncService(engine.GetDefaultContext())
    runner.RunService()
}
