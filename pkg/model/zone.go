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

// FindRecord finds a record.
func (z *Zone) FindRecord(name, typ string) (*ResourceRecordSet, bool) {
	for i, r := range z.Records {
		if r.Name == name && r.Type == typ {
			return &z.Records[i], true
		}
	}
	return nil, false
}
