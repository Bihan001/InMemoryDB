package config

var ServerHost string
var ServerPort int
var ConnectionLimit int

var KeyCountLimit int = 10
var EvictionMethod string = LRU_EVICTION
var EvictionPercentage float64 = 0.40

var LogFilePath = "./db.wal"

const (
    ORDERED_EVICTION = "ordered-eviction"
    RANDOM_EVICTION  = "random-eviction"
    LRU_EVICTION     = "lru-eviction"
)

const (
    EvictionPoolSize = 5
)