package engine

import (
	"github.com/Bihan001/MyDB/internal/ioprotocol"
	"github.com/Bihan001/MyDB/internal/stats"
	"github.com/Bihan001/MyDB/internal/store"
	"github.com/Bihan001/MyDB/internal/wal"
)

type Context struct {
    Encoder          ioprotocol.Encoder
    Decoder          ioprotocol.Decoder
    ExpiryManager    ExpiryManager
    StatsManager     stats.StatsManager
    EvictionManager  EvictionManager
    ObjectEncoder    ObjectEncoder
    WAL              wal.WAL
    Store            store.DataStore
    Evaluator        Evaluator
}

var defaultContext *Context

func init() {
    var statsManager stats.StatsManager = stats.GetNewStatsManager()
    var st store.DataStore = store.GetNewStore(statsManager)

    defaultContext = &Context{
        Encoder:         ioprotocol.GetNewRespEncoder(),
        Decoder:         ioprotocol.GetNewRespDecoder(),
        StatsManager:    statsManager,
        EvictionManager: GetNewEvictionManager(),
        ObjectEncoder:   GetNewObjectEncoder(),
        WAL:             wal.GetWAL(),
        Store:           st,
    }

    defaultContext.ExpiryManager = GetNewExpiryManager(defaultContext)
    defaultContext.Evaluator = GetNewEvaluator(defaultContext)
}

func GetDefaultContext() *Context {
    return defaultContext
}
