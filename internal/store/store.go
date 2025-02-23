package store

import (
    "time"

    "github.com/Bihan001/MyDB/internal/stats"
)

type DataStore interface {
    CreateEntry(content interface{}, expireInMs int64, objType uint8, objEnc uint8) *DataEntry
    Set(key string, entry *DataEntry)
    Get(key string) *DataEntry
    Del(key string) bool
    AllKeys() []string
    Size() int
}

type memoryStore struct {
    data         map[string]*DataEntry
    statsManager stats.StatsManager
}

type DataEntry struct {
    content      interface{}
    expiration   int64
    typeEncoding uint8
}

func GetNewStore(statsManager stats.StatsManager) DataStore {
    return &memoryStore{
        data: make(map[string]*DataEntry),
        statsManager: statsManager,
    }
}

func (de *DataEntry) GetValue() interface{} {
    return de.content
}

func (de *DataEntry) SetValue(v interface{}) {
    de.content = v
}

func (de *DataEntry) GetExpiration() int64 {
    return de.expiration
}

func (de *DataEntry) SetExpiration(e int64) {
    de.expiration = e
}

func (de *DataEntry) GetTypeEncoding() uint8 {
    return de.typeEncoding
}

func (ms *memoryStore) CreateEntry(content interface{}, expireInMs int64, objType uint8, objEnc uint8) *DataEntry {
    var exp int64 = -1
    if expireInMs > 0 {
        exp = time.Now().UnixMilli() + expireInMs
    }
    return &DataEntry{
        content:      content,
        expiration:   exp,
        typeEncoding: objType | objEnc,
    }
}

func (ms *memoryStore) Set(key string, entry *DataEntry) {
    ms.data[key] = entry
    ms.statsManager.IncrDBStat("keys")
}

func (ms *memoryStore) Get(key string) *DataEntry {
    val := ms.data[key]
    if val != nil {
        if val.expiration != -1 && val.expiration <= time.Now().UnixMilli() {
            ms.Del(key)
            return nil
        }
        return val
    }
    return nil
}

func (ms *memoryStore) Del(key string) bool {
    if _, exists := ms.data[key]; exists {
        delete(ms.data, key)
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
