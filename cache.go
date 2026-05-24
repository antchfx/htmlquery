package htmlquery

import (
	"errors"
	"sync"

	"github.com/antchfx/xpath"
	"github.com/golang/groupcache/lru"
)

type XpathQueryLookup interface {
	GetQuery(expr string) (*xpath.Expr, error)
}

var (
	// DisableSelectorCache will disable caching for the query selector if value is true.
	DisableSelectorCache = false

	// SelectorCacheMaxEntries allows how many selector object can be caching. Default is 50.
	// Will disable caching if SelectorCacheMaxEntries <= 0.
	SelectorCacheMaxEntries = 50
)

var (
	xpcache XpathQueryLookup
)

// max allows how many selector object can be caching. Default is 50.
// Will disable caching if max <= 0.
func NewXpathQueryLookup(max int) XpathQueryLookup {
	if max == 0 {
		return &nocacheXpathQueryLookup{}
	}
	return &lruXpathQueryLookup{
		cache: lru.New(max),
	}
}

type lruXpathQueryLookup struct {
	cache      *lru.Cache
	cacheMutex sync.Mutex
}

func (lxpl *lruXpathQueryLookup) GetQuery(expr string) (*xpath.Expr, error) {
	if lxpl.cache == nil || DisableSelectorCache {
		return xpath.Compile(expr)
	}

	lxpl.cacheMutex.Lock()
	defer lxpl.cacheMutex.Unlock()
	if v, ok := lxpl.cache.Get(expr); ok {
		e, ok := v.(*xpath.Expr)
		if !ok {
			return nil, errors.New("type asserion failed")
		}
		return e, nil
	}

	v, err := xpath.Compile(expr)
	if err != nil {
		return nil, err
	}
	lxpl.cache.Add(expr, v)
	return v, nil

}

type nocacheXpathQueryLookup struct{}

func (*nocacheXpathQueryLookup) GetQuery(expr string) (*xpath.Expr, error) {
	return xpath.Compile(expr)
}

func SetCache(x XpathQueryLookup) {
	xpcache = NewXpathQueryLookup(SelectorCacheMaxEntries)

}

func init() {
	SetCache(NewXpathQueryLookup(SelectorCacheMaxEntries))
}
