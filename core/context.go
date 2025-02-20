package core

import (
	ioprotocol "github.com/Bihan001/MyDB/core/io_protocol"
	"github.com/Bihan001/MyDB/core/stats"
	"github.com/Bihan001/MyDB/core/store"
)

type Context struct {
	Encoder ioprotocol.Encoder
	Decoder ioprotocol.Decoder
	ExpiryManager ExpiryManager
	StatsManager stats.StatsManager
	EvictionManager EvictionManager
	ObjectEncoder ObjectEncoder
	WAL WAL
	Store store.DataStore
}

var statsManager stats.StatsManager = stats.GetNewStatsManager();
var st store.DataStore = store.GetNewStore(statsManager);

var defaultContext *Context = &Context{
	Encoder: ioprotocol.GetNewRespEncoder(),
	Decoder: ioprotocol.GetNewRespDecoder(),
	ExpiryManager: GetNewExpiryManager(st),
	StatsManager: statsManager,
	EvictionManager: GetNewEvictionManager(),
	ObjectEncoder: GetNewObjectEncoder(),
	WAL: GetWAL(),
	Store: st,
}

func GetDefaultContext() *Context {
	return defaultContext
}
