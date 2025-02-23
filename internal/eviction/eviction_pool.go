package eviction

import (
	"sort"

	"github.com/Bihan001/MyDB/internal/config"
	"github.com/Bihan001/MyDB/internal/utils"
)

type PoolElement struct {
	key string
	lastAccessedAtMs uint64
}

type EvictionPool struct {
	elements []*PoolElement
	keySet map[string]bool
}

func GetNewEvictionPool() *EvictionPool {
	return &EvictionPool{
		elements: make([]*PoolElement, 0),
		keySet: make(map[string]bool),
	}
}

type SortHelper []*PoolElement

func (sh SortHelper) Len() int {
	return len(sh)
}

func (sh SortHelper) Less(i, j int) bool {
	return sh[i].lastAccessedAtMs < sh[j].lastAccessedAtMs
}

func (sh SortHelper) Swap(i, j int) {
	sh[i], sh[j] = sh[j], sh[i]
}

func (ep *EvictionPool) Push(key string, lastAccessedAtMs uint64) {

	if _, exists := ep.keySet[key]; exists {
		return
	}

	if len(ep.elements) < config.EvictionPoolSize {
		ep.elements = append(ep.elements, &PoolElement{key: key, lastAccessedAtMs: lastAccessedAtMs})
		ep.keySet[key] = true
		sort.Sort(SortHelper(ep.elements))
	} else if lastAccessedAtMs < ep.elements[0].lastAccessedAtMs {
		ep.keySet[ep.elements[0].key] = false
		ep.elements[0] = &PoolElement{key: key, lastAccessedAtMs: lastAccessedAtMs}
		ep.keySet[key] = true
	}

}

func (ep *EvictionPool) Pop() *PoolElement {
	if len(ep.elements) == 0 {
		return nil
	}
	el := ep.elements[0]
	ep.elements = ep.elements[1:]
	delete(ep.keySet, el.key)
	return el
}

func (ep *EvictionPool) Size() int {
	return len(ep.elements)
}

func (ep *EvictionPool) UpdateLastAccessedTime(key string) {
	if _, exists := ep.keySet[key]; !exists {
		return
	}

	for _, el := range ep.elements {
		if el.key == key {
			el.lastAccessedAtMs = utils.GetCurrentTime()
			sort.Sort(SortHelper(ep.elements))
			return
		}
	}
}