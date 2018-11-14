package cmd

import (
	"encoding/json"
	"os"

	"github.com/bitnami-labs/zerozone/pkg/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	fileCfg = "file"
)

var (
	recordName string
	recordType string
	recordData string
)

func registerEditFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(fileCfg, "f", "", "Zone file")
	cmd.Flags().StringVarP(&recordName, "name", "n", "", "Record name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&recordType, "type", "t", "", "Record type")
	cmd.MarkFlagRequired("type")
	cmd.Flags().StringVarP(&recordData, "data", "d", "", "record data")

	viper.BindPFlag(fileCfg, cmd.Flags().Lookup(fileCfg))
}

func openZone(fileName string) (zone *model.Zone, save func() error, err error) {
	save = func() error {
		f, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer f.Close()
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		return enc.Encode(zone)
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
