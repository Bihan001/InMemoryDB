package eviction

import (
	"github.com/Bihan001/MyDB/internal/config"
	"github.com/Bihan001/MyDB/internal/interfaces"
)

type EvictionPolicy interface {
	Evict()
}

type OrderedEvictionPolicy struct {
	store interfaces.DataStore
}

type RandomEvictionPolicy struct {
	store interfaces.DataStore
}

type LRUEvictionPolicy struct {
	store        interfaces.DataStore
	evictionPool *EvictionPool
}

func GetNewEvictionPolicy(store interfaces.DataStore, evictionPool *EvictionPool) EvictionPolicy {
	switch config.EvictionMethod {
	case config.ORDERED_EVICTION:
		return &OrderedEvictionPolicy{
			store: store,
		}
	case config.RANDOM_EVICTION:
		return &RandomEvictionPolicy{
			store: store,
		}
	case config.LRU_EVICTION:
		return &LRUEvictionPolicy{
			store:        store,
			evictionPool: evictionPool,
		}
	default:
		return nil
	}
}

func (ep *OrderedEvictionPolicy) Evict() {
	keys := ep.store.AllKeys()
	if len(keys) == 0 {
		return
	}
	ep.store.Del(keys[0])
}

func (ep *RandomEvictionPolicy) Evict() {
	keys := ep.store.AllKeys()
	if len(keys) == 0 {
		return
	}
	toRemove := int64(config.EvictionPercentage * float64(config.KeyCountLimit))
	idx := 0
	for idx < len(keys) && toRemove > 0 {
		ep.store.Del(keys[idx])
		toRemove--
		idx++
	}
}

func (ep *LRUEvictionPolicy) Evict() {
	populateLruEvictionPool(ep.store, ep.evictionPool)

	evictionCount := int(config.EvictionPercentage * float64(config.KeyCountLimit))

	for i := 0; i < evictionCount && ep.evictionPool.Size() > 0; i++ {
		elem := ep.evictionPool.Pop()
		if elem == nil {
			return
		}
		ep.store.Del(elem.key)
	}
}

func populateLruEvictionPool(st interfaces.DataStore, evictionPool *EvictionPool) {
	n := 5

	for _, key := range st.AllKeys() {
		entry := st.Get(key)
		evictionPool.Push(key, entry.GetLastAccessedMs())
		n--
		if n == 0 {
			break
		}
	}
}
