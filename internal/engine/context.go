package engine

import (
	"github.com/Bihan001/MyDB/internal/eviction"
	"github.com/Bihan001/MyDB/internal/expiration"
	"github.com/Bihan001/MyDB/internal/interfaces"
	"github.com/Bihan001/MyDB/internal/ioprotocol"
	"github.com/Bihan001/MyDB/internal/object_encoding"
	"github.com/Bihan001/MyDB/internal/stats"
	"github.com/Bihan001/MyDB/internal/store"
	"github.com/Bihan001/MyDB/internal/wal"
)

var defaultContext *interfaces.Context

func init() {
	statsManager := stats.GetNewStatsManager()
	evictionPool := eviction.GetNewEvictionPool()
	store := store.GetNewStore(statsManager, evictionPool.UpdateLastAccessedTime)

	defaultContext = &interfaces.Context{
		Encoder:        ioprotocol.GetNewRespEncoder(),
		Decoder:        ioprotocol.GetNewRespDecoder(),
		StatsManager:   statsManager,
		EvictionPolicy: eviction.GetNewEvictionPolicy(store, evictionPool),
		ObjectEncoder:  object_encoding.GetNewObjectEncoder(),
		WAL:            wal.GetWAL(),
		Store:          store,
	}

	defaultContext.ExpiryManager = expiration.GetNewExpiryManager(defaultContext)
	defaultContext.Evaluator = GetNewEvaluator(defaultContext)
}

func GetDefaultContext() *interfaces.Context {
	return defaultContext
}
