package cmd

import (
	"fmt"
	"os"

	"github.com/bitnami-labs/zerozone/pkg/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	recordTTL uint32
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new DNS record",
	Long:  `Add a new DNS record`,
	RunE: func(cmd *cobra.Command, args []string) error {
		zone, save, err := openZone(viper.GetString(fileCfg))
		if err != nil {
			return err
		}
		r, ok := zone.FindRecord(recordName, recordType)
		if !ok {
			zone.Records = append(zone.Records, model.ResourceRecordSet{
				Name: recordName,
				Type: recordType,
			})
			r = &zone.Records[len(zone.Records)-1]
		}
		r.TTL = recordTTL

		found := false
		for _, d := range r.RRDatas {
			if d == recordData {
				found = true
			}
		}
		if !found {
			r.RRDatas = append(r.RRDatas, recordData)
		}

		zoneName, err := getZoneName()
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Printf("%s.%s.%s\n", recordName, zoneName, viper.GetString(zeroZoneDomainCfg))
		fmt.Fprintf(os.Stderr, "\n")

		return save()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	registerEditFlags(addCmd)

	addCmd.Flags().Uint32Var(&recordTTL, "ttl", 60, "record ttl")
	addCmd.MarkFlagRequired("data")

	viper.BindPFlag("file", addCmd.Flags().Lookup("file"))
}
