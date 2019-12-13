package htmlquery

import (
	"sync"

	"github.com/antchfx/xpath"
	"github.com/golang/groupcache/lru"
)

// DisableSelectorCache will disable caching for the query selector if value is true.
var DisableSelectorCache = false

// SelectorCacheMaxEntries allows how many selector object can be caching. Default is 50.
// Will disable caching if SelectorCacheMaxEntries <= 0.
var SelectorCacheMaxEntries = 50

var (
	cacheOnce  sync.Once
	cache      *lru.Cache
	cacheMutex sync.RWMutex
)

func getQuery(expr string) (*xpath.Expr, error) {
	if DisableSelectorCache || SelectorCacheMaxEntries <= 0 {
		return xpath.Compile(expr)
	}
	cacheOnce.Do(func() {
		cache = lru.New(50)
	})
	cacheMutex.RLock()
	if v, ok := cache.Get(expr); ok {
		cacheMutex.RUnlock()
		return v.(*xpath.Expr), nil
	}
	cacheMutex.RUnlock()
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	v, err := xpath.Compile(expr)
	if err != nil {
		return nil, err
	}
	cache.Add(expr, v)
	return v, nil

}
