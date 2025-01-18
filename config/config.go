package config

var Host string
var Port int
var MaxClients int

var MaxKeyLimit int = 5
var EvictionStrategy string = EVICT_SIMPLE_FIRST

var WALFilePath = "./db.wal"

const (
	EVICT_SIMPLE_FIRST = "evict-simple-first"
)