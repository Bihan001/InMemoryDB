package engine

import (
    "github.com/Bihan001/MyDB/internal/config"
    "github.com/Bihan001/MyDB/internal/store"
)

type OrderedEvictionPolicy struct {}

type RandomEvictionPolicy struct {}

func (oep *OrderedEvictionPolicy) Evict(st store.DataStore) {
    keys := st.AllKeys()
    if len(keys) == 0 {
        return
    }
    st.Del(keys[0])
}

func (rep *RandomEvictionPolicy) Evict(st store.DataStore) {
    keys := st.AllKeys()
    if len(keys) == 0 {
        return
    }
    toRemove := int64(config.EvictionPercentage * float64(config.KeyCountLimit))
    idx := 0
    for idx < len(keys) && toRemove > 0 {
        st.Del(keys[idx])
        toRemove--
        idx++
    }
}
