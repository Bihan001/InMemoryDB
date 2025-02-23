package engine

import (
	"github.com/Bihan001/MyDB/internal/eviction"
	"github.com/Bihan001/MyDB/internal/interfaces"
	"github.com/Bihan001/MyDB/internal/ioprotocol"
	"github.com/Bihan001/MyDB/internal/stats"
	"github.com/Bihan001/MyDB/internal/store"
	"github.com/Bihan001/MyDB/internal/wal"
)

type Context struct {
	Encoder        ioprotocol.Encoder
	Decoder        ioprotocol.Decoder
	ExpiryManager  ExpiryManager
	StatsManager   stats.StatsManager
	EvictionPolicy eviction.EvictionPolicy
	ObjectEncoder  ObjectEncoder
	WAL            wal.WAL
	Store          interfaces.DataStore
	Evaluator      Evaluator
}

var defaultContext *Context

func init() {
	var statsManager stats.StatsManager = stats.GetNewStatsManager()
	evictionPool := eviction.GetNewEvictionPool()
	var st interfaces.DataStore = store.GetNewStore(statsManager, evictionPool.UpdateLastAccessedTime)

	defaultContext = &Context{
		Encoder:        ioprotocol.GetNewRespEncoder(),
		Decoder:        ioprotocol.GetNewRespDecoder(),
		StatsManager:   statsManager,
		EvictionPolicy: eviction.GetNewEvictionPolicy(st, evictionPool),
		ObjectEncoder:  GetNewObjectEncoder(),
		WAL:            wal.GetWAL(),
		Store:          st,
	}

	defaultContext.ExpiryManager = GetNewExpiryManager(defaultContext)
	defaultContext.Evaluator = GetNewEvaluator(defaultContext)
}

func GetDefaultContext() *Context {
	return defaultContext
}
