package core

import "time"

type Store struct {
	store map[string]*Value
}

type Value struct {
	value interface{}
	expiresAt int64
}

var store *Store

func init() {
	if store == nil {
		store = &Store{
			store: make(map[string]*Value),
		}
	}
}

func (s *Store) NewValue(value interface{}, expiryDurationMs int64) *Value {
	var expiresAt int64 = -1

	if expiryDurationMs > 0 {
		expiresAt = time.Now().UnixMilli() + expiryDurationMs
	}

	return &Value{
		value: value,
		expiresAt: expiresAt,
	}
}

func (s *Store) Put(k string, v *Value) {
	s.store[k] = v
}

func (s *Store) Get(k string) *Value {
	v := s.store[k]

	if v != nil {
		if v.expiresAt != -1 && v.expiresAt <= time.Now().UnixMilli() {
			delete(s.store, k)
			return nil
		}
		return v
	}
	return nil
}

func (s *Store) Delete(k string) bool {
	if _, ok := s.store[k]; ok {
		delete(s.store, k)
		return true
	}
	return false
}
