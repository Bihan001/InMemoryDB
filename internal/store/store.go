package store

import (
	"time"

	"github.com/Bihan001/MyDB/internal/interfaces"
	"github.com/Bihan001/MyDB/internal/stats"
	"github.com/Bihan001/MyDB/internal/utils"
)

type OnKeyAccess func(key string)

type memoryStore struct {
    data         map[string]*interfaces.DataEntry
    expiryMap map[*interfaces.DataEntry]uint64
    statsManager stats.StatsManager
    onKeyAccess OnKeyAccess
}

func GetNewStore(statsManager stats.StatsManager, onKeyAccess OnKeyAccess) interfaces.DataStore {
    return &memoryStore{
        data: make(map[string]*interfaces.DataEntry),
        expiryMap: make(map[*interfaces.DataEntry]uint64),
        statsManager: statsManager,
        onKeyAccess: onKeyAccess,
    }
}

func (ms *memoryStore) CreateEntry(content interface{}, expireInMs int64, objType uint8, objEnc uint8) *interfaces.DataEntry {
    entry := interfaces.CreateDataEntry(content, objType | objEnc, utils.GetCurrentTime())

    if expireInMs > 0 {
        ms.SetExpiry(entry, expireInMs)
    }

    return entry
}

func (ms *memoryStore) HasExpired(entry *interfaces.DataEntry) bool {
    if expiry, exists := ms.expiryMap[entry]; exists {
        return expiry <= uint64(time.Now().UnixMilli())
    }
    return false
}

func (ms *memoryStore) GetExpiryMs(entry *interfaces.DataEntry) (uint64, bool) {
    expiry, exists := ms.expiryMap[entry]
    return expiry, exists
}

func (ms *memoryStore) SetExpiry(entry *interfaces.DataEntry, expireInMs int64) {
    ms.expiryMap[entry] = uint64(time.Now().UnixMilli()) + uint64(expireInMs)
}

func (ms *memoryStore) Set(key string, entry *interfaces.DataEntry) {
    ms.setLastAccessedTime(key, entry)
    ms.data[key] = entry
    ms.statsManager.IncrDBStat("keys")
}

func (ms *memoryStore) Get(key string) *interfaces.DataEntry {
    val := ms.data[key]
    if val != nil {
        if ms.HasExpired(val) {
            ms.Del(key)
            return nil
        }
        ms.setLastAccessedTime(key, val)
        return val
    }
    return nil
}

func (ms *memoryStore) Del(key string) bool {
    if val, exists := ms.data[key]; exists {
        delete(ms.data, key)
        delete(ms.expiryMap, val)
        ms.statsManager.DecrDBStat("keys")
        return true
    }
    return false
}

func (ms *memoryStore) AllKeys() []string {
    keys := make([]string, 0, len(ms.data))
    for k := range ms.data {
        keys = append(keys, k)
    }
    return keys
}

func (ms *memoryStore) Size() int {
    return len(ms.data)
}

func (ms *memoryStore) setLastAccessedTime(key string, entry *interfaces.DataEntry) {
    entry.SetLastAccessedMs(utils.GetCurrentTime())
    ms.onKeyAccess(key)
}