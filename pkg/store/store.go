// Package store abstracts access to zero config zones.
package store

import (
	"github.com/bitnami-labs/zerozone/pkg/model"
)

// A Fetcher knows how to fetch a zone.
type Fetcher interface {
	FetchZone(id string) (*model.Zone, error)
}
