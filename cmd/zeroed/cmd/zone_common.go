package cmd

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/multiformats/go-multibase"
	"github.com/spf13/viper"
)

type KeyListResult struct {
	Keys []struct {
		Name string
		Id   string
	}
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