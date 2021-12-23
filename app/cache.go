package app

import (
	"sort"
	"sync"
)

//qCache -
type queueCache struct {
	items    sort.StringSlice
	mu       sync.Mutex
	unsorted bool
}

var qCache = queueCache{items: make(sort.StringSlice, 0)}

func addCache(item string) {
	qCache.mu.Lock()
	qCache.items = append(qCache.items, item)
	qCache.unsorted = true
	qCache.mu.Unlock()
}
