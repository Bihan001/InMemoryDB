package interfaces

type DataStore interface {
    CreateEntry(content interface{}, expireInMs int64, objType uint8, objEnc uint8) *DataEntry
    Set(key string, entry *DataEntry)
    Get(key string) *DataEntry
    Del(key string) bool
    AllKeys() []string
    Size() int
}

type Expirable interface {
	HasExpired(entry *DataEntry) bool
	GetExpiryMs(entry *DataEntry) (uint64, bool)
	SetExpiry(entry *DataEntry, expireInMs int64)
}

type DataEntry struct {
    content      interface{}
    lastAccessedAtMs   uint64
    typeEncoding uint8
}

func CreateDataEntry(content interface{}, typeEncoding uint8, lastAccessAtMs uint64) *DataEntry {
	return &DataEntry{
		content: content,
		lastAccessedAtMs: lastAccessAtMs,
		typeEncoding: typeEncoding,
	}
}

func (de *DataEntry) GetValue() interface{} {
    return de.content
}

func (de *DataEntry) SetValue(v interface{}) {
    de.content = v
}

func (de *DataEntry) GetLastAccessedMs() uint64 {
    return de.lastAccessedAtMs
}

func (de *DataEntry) SetLastAccessedMs(ts uint64) {
    de.lastAccessedAtMs = ts
}

func (de *DataEntry) GetTypeEncoding() uint8 {
    return de.typeEncoding
}