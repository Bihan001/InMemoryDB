package interfaces

type ObjectEncoder interface {
    EvaluateObjectEncoding(val string) (uint8, uint8)
}

type Evaluator interface {
    Evaluate(ops OperationList) ([]byte, error)
}

type EvictionPolicy interface {
	Evict()
}

type ExpiryManager interface {
    PurgeExpiredEntries()
}

type StatsManager interface {
    RecordDBStat(metric string, value int)
    IncrDBStat(metric string)
    DecrDBStat(metric string)
    GetDbStats() [4]map[string]int
}


