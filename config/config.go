package config

var Host string
var Port int
var MaxClients int

var MaxKeyLimit int = 5
var EvictionStrategy string = EVICT_SIMPLE_FIRST
var EvictionRatio float64 = 0.40

var WALFilePath = "./db.wal"

const (
	EVICT_SIMPLE_FIRST = "evict-simple-first"
	EVICT_ALL_KEYS_RANDOM = "all-keys-random"
)