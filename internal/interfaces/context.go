package interfaces

import (
	"github.com/Bihan001/MyDB/internal/ioprotocol"
	"github.com/Bihan001/MyDB/internal/wal"
)

type Context struct {
	Encoder        ioprotocol.Encoder
	Decoder        ioprotocol.Decoder
	ExpiryManager  ExpiryManager
	StatsManager   StatsManager
	EvictionPolicy EvictionPolicy
	ObjectEncoder  ObjectEncoder
	WAL            wal.WAL
	Store          DataStore
	Evaluator      Evaluator
}
