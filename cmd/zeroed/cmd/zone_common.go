package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bitnami-labs/zerozone/pkg/model"
	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/multiformats/go-multibase"
	"github.com/spf13/viper"
)

// KeyListResult is the result of the key/list IPFS API method.
// It's here because it's not exported by go-ipfs-api.
type KeyListResult struct {
	Keys []struct {
		Name string
		Id   string
	}
}

func openZone(fileName string) (zone *model.Zone, save func() error, err error) {
	save = func() error {
		return updateZone(fileName, zone)
	}

	f, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return &model.Zone{}, save, nil
	}
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&zone); err != nil {
		return nil, nil, err
	}
	return zone, save, nil
}

func updateZone(fileName string, zone *model.Zone) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	zone.Generation++
	if err := enc.Encode(zone); err != nil {
		return err
	}

	if noPublish {
		return nil
	}

	f.Seek(0, 0)
	return publishZone(f)
}

func publishZone(r io.Reader) error {
	sh := shell.NewShell(viper.GetString(apiAddrCfg))
	hash, err := sh.Add(r)
	if err != nil {
		return err
	}
	pubkey := viper.GetString(pubKeyCfg)
	fmt.Fprintf(os.Stderr, "Publishing /ipfs/%s to IPNS\n", hash)
	_, err = sh.PublishWithDetails(hash, pubkey, 7*24*time.Hour, 30*time.Second, false)
	return err
}

func getZoneName() (string, error) {
	pubKey, err := getPubKey()
	if err != nil {
		return "", err
	}
	c0, err := cid.Decode(pubKey)
	if err != nil {
		return "", err
	}
	c := cid.NewCidV1(cid.DagProtobuf, c0.Hash())
	return c.Encode(multibase.MustNewEncoder(multibase.Base32)), nil
}

func getPubKey() (string, error) {
	sh := shell.NewShell(viper.GetString(apiAddrCfg))
	var out KeyListResult
	if err := sh.Request("key/list").Exec(context.Background(), &out); err != nil {
		return "", err
	}
	pubKeyName := viper.GetString(pubKeyCfg)
	for _, k := range out.Keys {
		if k.Name == pubKeyName {
			return k.Id, nil
		}
	}
	return "", fmt.Errorf("cannot find pubkey %q", pubKeyName)
}