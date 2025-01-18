package core

import "github.com/Bihan001/MyDB/config"

func evict() {
	switch config.EvictionStrategy {
	case config.EVICT_SIMPLE_FIRST:
		evictFirst()
	}
}

func evictFirst() {
	for key := range store.store {
		delete(store.store, key)
		return
	}
}
