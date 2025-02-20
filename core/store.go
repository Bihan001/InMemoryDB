package core

import (
	"time"

	"github.com/Bihan001/MyDB/config"
)

type Store struct {
	store map[string]*Value
}

type Value struct {
	value interface{}
	expiresAt int64
	typeEncoding uint8
}

var OBJECT_TYPE_STRING uint8 = 0

var OBJ_ENCODING_INT uint8 = 1
var OBJ_ENCODING_RAW uint8 = 0
var OBJ_ENCODING_EMBEDDED_STR uint8 = 8

var store *Store

func init() {
	if store == nil {
		store = &Store{
			store: make(map[string]*Value),
		}
	}
}

func (s *Store) NewValue(value interface{}, expiryDurationMs int64, objectType uint8, objectEncoding uint8) *Value {
	var expiresAt int64 = -1

	if expiryDurationMs > 0 {
		expiresAt = time.Now().UnixMilli() + expiryDurationMs
	}

	return &Value{
		value: value,
		expiresAt: expiresAt,
		typeEncoding: objectType|objectEncoding,
	}
}

func (s *Store) Put(k string, v *Value) {
	if len(s.store) >= config.MaxKeyLimit {
		evict()
	}
	s.store[k] = v
	if KeyspaceStat[0] == nil {
		KeyspaceStat[0] = make(map[string]int)
	}
	KeyspaceStat[0]["keys"]++
}

func (s *Store) Get(k string) *Value {
	v := s.store[k]

	if v != nil {
		if v.expiresAt != -1 && v.expiresAt <= time.Now().UnixMilli() {
			s.Delete(k)
			return nil
		}
		return v
	}
	return nil
}

func (s *Store) Delete(k string) bool {
	if _, ok := s.store[k]; ok {
		delete(s.store, k)
		KeyspaceStat[0]["keys"]--
		return true
	}
	return false
}

func (s *Store) Keys() []string {
	keys := make([]string, len(s.store))
	i := 0
	for k := range s.store {
		keys[i] = k
		i++
	}
	return keys
}

func (s *Store) Length() int {
	return len(s.store)
}