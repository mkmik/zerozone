package store

import (
	"sync"

	"github.com/bitnami-labs/zerozone/pkg/model"
	"github.com/coredns/coredns/plugin/pkg/log"
)

// CachingFetcher is a simple stupid in memory cache that unconditionally returns
// cached zones and triggers a cache update in the background.
type CachingFetcher struct {
	fetcher Fetcher
	cache   sync.Map
}

// NewCachingFetcher creates a new caching fetcher
func NewCachingFetcher(fetcher Fetcher) *CachingFetcher {
	return &CachingFetcher{fetcher: fetcher}
}

func (f *CachingFetcher) FetchZone(id string) (*model.Zone, error) {
	zi, ok := f.cache.Load(id)
	if ok {
		log.Debugf("returning cached entry for %q, triggering cache update in the background", id)
		go f.refresh(id)
		return zi.(*model.Zone), nil
	}

	return f.savingFetchZone(id)
}

func (f *CachingFetcher) refresh(id string) {
	if _, err := f.savingFetchZone(id); err != nil {
		log.Errorf("failed to refresh cache for %q: %v", id, err)
	}
}

func (f *CachingFetcher) savingFetchZone(id string) (*model.Zone, error) {
	z, err := f.fetcher.FetchZone(id)
	if err != nil {
		return nil, err
	}

	log.Debugf("storing %q in cache", id)
	f.cache.Store(id, z)
	return z, nil
}
