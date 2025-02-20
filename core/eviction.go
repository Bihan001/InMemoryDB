package core

import "github.com/Bihan001/MyDB/config"

func evict() {
	switch config.EvictionStrategy {
	case config.EVICT_SIMPLE_FIRST:
		evictFirst()
	case config.EVICT_ALL_KEYS_RANDOM:
		evictAllkeysRandom()
	}
}

func evictFirst() {
	for key := range store.store {
		store.Delete(key)
		return
	}
}

func evictAllkeysRandom() {
	evictCount := int64(config.EvictionRatio * float64(config.MaxKeyLimit))
	for k := range store.store {
		store.Delete(k)
		evictCount--
		if evictCount <= 0 {
			break
		}
	}
}