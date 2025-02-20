package core

import (
	"github.com/Bihan001/MyDB/config"
	"github.com/Bihan001/MyDB/core/store"
)

type OrderedEvictionPolicy struct {}

type RandomEvictionPolicy struct {}

func (oep *OrderedEvictionPolicy) Evict(st store.DataStore) {
    keys := st.AllKeys()
    if len(keys) == 0 {
        return
    }
    st.Remove(keys[0])
}

func (rep *RandomEvictionPolicy) Evict(st store.DataStore) {
    keys := st.AllKeys()
    if len(keys) == 0 {
        return
    }
    toRemove := int64(config.EvictionPercentage * float64(config.KeyCountLimit))
    idx := 0
    for idx < len(keys) && toRemove > 0 {
        st.Remove(keys[idx])
        toRemove--
        idx++
    }
}
