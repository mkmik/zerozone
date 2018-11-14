package cmd

import (
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
