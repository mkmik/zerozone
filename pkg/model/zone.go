// Package model defines the data model and our JSON encoding of Zero Zones.
package model

// Zone is a user's zone.
type Zone struct {
	Records []ResourceRecordSet `json:"records,omitempty"`
}

// ResourceRecordSet is a unit of data that will be returned by the DNS servers.
type ResourceRecordSet struct {
	Name    string   `json:"name,omitempty"`
	Type    string   `json:"type,omitempty"`
	TTL     uint32   `json:"ttl,omitempty"`
	RRDatas []string `json:"rrdatas,omitempty"`
}
