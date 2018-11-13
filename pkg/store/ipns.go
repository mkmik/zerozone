package store

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitnami-labs/zerozone/pkg/model"
	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	multibase "github.com/multiformats/go-multibase"
)

// IPNSFetcher knows how to fetch zones from IPNS.
type IPNSFetcher struct {
	shell *shell.Shell
}

// NewIPNSFetcher returns a new IPNSFetcher.
func NewIPNSFetcher(apiAddr string) *IPNSFetcher {
	return &IPNSFetcher{
		shell: shell.NewShell(apiAddr),
	}
}

func (f *IPNSFetcher) FetchZone(id string) (*model.Zone, error) {
	zoneAddr, err := ipnsAddr(id)
	if err != nil {
		return nil, err
	}

	rs, err := f.shell.Cat(zoneAddr)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var zone model.Zone
	if err := json.NewDecoder(rs).Decode(&zone); err != nil {
		return nil, err
	}

	return &zone, nil
}

// IPNSGatewayFetcher fetches a zone using the ipfs gateway
type IPNSGatewayFetcher struct {
	gwAddr string
}

// NewIPNSGatewayFetcher returns a new IPNSGatewayFetcher.
func NewIPNSGatewayFetcher(gwAddr string) *IPNSGatewayFetcher {
	return &IPNSGatewayFetcher{
		gwAddr: gwAddr,
	}
}

func (f *IPNSGatewayFetcher) FetchZone(id string) (*model.Zone, error) {
	zoneAddr, err := ipnsAddr(id)
	if err != nil {
		return nil, err
	}

	httpAddr := fmt.Sprintf("%s%s", f.gwAddr, zoneAddr)
	log.Debugf("fetching from gw: %q", httpAddr)
	rs, err := http.Get(httpAddr)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()

	var zone model.Zone
	if err := json.NewDecoder(rs.Body).Decode(&zone); err != nil {
		return nil, err
	}

	return &zone, nil
}

func ipnsAddr(hash string) (string, error) {
	// ipns addresses cannot yet be V1 cid addresses.
	legacy, err := toLegacyBase58(hash)
	if err != nil {
		return "", err
	}

	addr := fmt.Sprintf("/ipns/%s", legacy)
	log.Debugf("addr %s", addr)
	return addr, nil
}

func toLegacyBase58(hash string) (string, error) {
	log.Debugf("parsing cid %q", hash)
	v1id, err := cid.Decode(hash)
	if err != nil {
		return "", err
	}
	v0id := cid.NewCidV0(v1id.Hash())
	return v0id.Encode(multibase.MustNewEncoder(multibase.Base58BTC)), nil
}
