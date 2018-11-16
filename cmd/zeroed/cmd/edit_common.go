package cmd

import (
	"github.com/spf13/cobra"
)

const (
	fileCfg = "file"
)

var (
	recordName string
	recordType string
	recordData string

	noPublish bool
)

func registerEditFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(fileCfg, "f", "", "Zone file")
	cmd.Flags().StringVarP(&recordName, "name", "n", "", "Record name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&recordType, "type", "t", "", "Record type")
	cmd.MarkFlagRequired("type")
	cmd.Flags().StringVarP(&recordData, "data", "d", "", "record data")

	cmd.Flags().BoolVar(&noPublish, "no-publish", false, "skip publishing to IPNS, but update local zone file")
}
