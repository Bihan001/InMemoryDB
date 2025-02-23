package engine

import (
	"log"

	"github.com/Bihan001/MyDB/internal/interfaces"
)

type ExpiryManager interface {
    PurgeExpiredEntries()
}

type defaultExpiryManager struct {
    sampleSize int
    context    *Context
    expirableStore interfaces.Expirable
}

func GetNewExpiryManager(context *Context) ExpiryManager {
    return &defaultExpiryManager{
        sampleSize: 20,
        context:    context,
        expirableStore: context.Store.(interfaces.Expirable),
    }
}

func (eas *defaultExpiryManager) PurgeExpiredEntries() {
    for {
        fraction := eas.scanAndRemoveExpired()
        if fraction < 0.25 {
            break
        }
    }
    log.Println("purged expired keys. total keys", eas.context.Store.Size())
}

func (eas *defaultExpiryManager) scanAndRemoveExpired() float32 {
    localSample := eas.sampleSize
    removed := 0
    keys := eas.context.Store.AllKeys()

    for _, k := range keys {
        if localSample == 0 {
            break
        }
        localSample--
        entry := eas.context.Store.Get(k)
        if entry == nil {
            continue
        }
        if eas.expirableStore.HasExpired(entry) {
            eas.context.Store.Del(k)
            removed++
        }
    }
    return float32(removed) / float32(eas.sampleSize)
}
