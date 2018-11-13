package store

import (
	"github.com/bitnami-labs/zerozone/pkg/model"
	"github.com/coredns/coredns/plugin/pkg/log"
	"golang.org/x/sync/singleflight"
)

type SingleFlightFetcher struct {
	fetcher Fetcher
	group   singleflight.Group
}

// NewSingleFlightFetcher creates a new single flight fetcher.
func NewSingleFlightFetcher(fetcher Fetcher) *SingleFlightFetcher {
	return &SingleFlightFetcher{fetcher: fetcher}
}

func (f *SingleFlightFetcher) FetchZone(id string) (*model.Zone, error) {
	v, err, _ := f.group.Do(id, func() (interface{}, error) {
		log.Debugf("fetching %q", id)
		return f.fetcher.FetchZone(id)
	})
	if err != nil {
		return nil, err
	}
	return v.(*model.Zone), nil
}
