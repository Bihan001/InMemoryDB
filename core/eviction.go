package core

import (
	"github.com/Bihan001/MyDB/config"
	"github.com/Bihan001/MyDB/core/store"
)

type EvictionManager interface {
    Evict(st store.DataStore)
}

type defaultEvictionManager struct {}

func GetNewEvictionManager() EvictionManager {
    return &defaultEvictionManager{}
}

func (eo *defaultEvictionManager) Evict(st store.DataStore) {
    switch config.EvictionMethod {
    case config.ORDERED_EVICTION:
        pol := &OrderedEvictionPolicy{}
        pol.Evict(st)
    case config.RANDOM_EVICTION:
        pol := &RandomEvictionPolicy{}
        pol.Evict(st)
    }
}
