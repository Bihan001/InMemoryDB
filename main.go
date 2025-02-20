package main

import (
	"flag"

	"github.com/Bihan001/MyDB/config"
	"github.com/Bihan001/MyDB/core"
	"github.com/Bihan001/MyDB/server"
)

func setupFlags() {
    flag.StringVar(&config.ServerHost, "host", "0.0.0.0", "Server Host")
    flag.IntVar(&config.ServerPort, "port", 6389, "Server Port")
    flag.IntVar(&config.ConnectionLimit, "maxClients", 20000, "Max Clients to support")
    flag.Parse()
}

func main() {
    setupFlags()

    var evaluator core.Evaluator = core.GetNewEvaluator(core.GetDefaultContext())
    var runner server.ServiceRunner = server.NewAsyncService(core.GetDefaultContext(), evaluator)
    
    runner.RunService()
}
