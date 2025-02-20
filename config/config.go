package config

var ServerHost string
var ServerPort int
var ConnectionLimit int

var KeyCountLimit int = 5
var EvictionMethod string = ORDERED_EVICTION
var EvictionPercentage float64 = 0.40

var LogFilePath = "./db.wal"

const (
    ORDERED_EVICTION = "ordered-eviction"
    RANDOM_EVICTION  = "random-eviction"
)
