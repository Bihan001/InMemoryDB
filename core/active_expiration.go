package core

import (
	"log"
	"time"

	"github.com/Bihan001/MyDB/core/store"
)

type ExpiryManager interface {
    PurgeExpiredEntries()
}

type defaultExpiryManager struct {
    sampleSize int
    store      store.DataStore
}

func GetNewExpiryManager(st store.DataStore) ExpiryManager {
    return &defaultExpiryManager{
        sampleSize: 20,
        store:      st,
    }
}

func (eas *defaultExpiryManager) PurgeExpiredEntries() {
    for {
        fraction := eas.scanAndRemoveExpired()
        if fraction < 0.25 {
            break
        }
    }
    log.Println("purged expired keys. total keys", eas.store.Size())
}

func (eas *defaultExpiryManager) scanAndRemoveExpired() float32 {
    localSample := eas.sampleSize
    removed := 0
    keys := eas.store.AllKeys()

    for _, k := range keys {
        if localSample == 0 {
            break
        }
        localSample--
        entry := eas.store.Retrieve(k)
        if entry == nil {
            continue
        }
        if entry.GetExpiration() != -1 && entry.GetExpiration() <= time.Now().UnixMilli() {
            eas.store.Remove(k)
            removed++
        }
    }
    return float32(removed) / float32(eas.sampleSize)
}

