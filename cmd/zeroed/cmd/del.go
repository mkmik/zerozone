package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete a resource record",
	Long:  `Delete a resource record from a Zero Zone`,
	RunE: func(cmd *cobra.Command, args []string) error {
		zone, save, err := openZone(viper.GetString(fileCfg))
		if err != nil {
			return err
		}
		r, ok := zone.FindRecord(recordName, recordType)
		if !ok {
			return fmt.Errorf("cannot find %q %q", recordName, recordType)
		}
		if recordData == "" {
			r.RRDatas = nil
		} else {
			for i := range r.RRDatas {
				found := false
				if r.RRDatas[i] == recordData {
					r.RRDatas = append(r.RRDatas[:i], r.RRDatas[i+1:]...)
					found = true
					break
				}
				if !found {
					return fmt.Errorf("cannot find %q %q %q", recordName, recordType, recordData)
				}
			}
		}
		if len(r.RRDatas) == 0 {
			for i := range zone.Records {
				if &zone.Records[i] == r {
					zone.Records = append(zone.Records[:i], zone.Records[i+1:]...)
					break
				}
			}
		}

		return save()
	},
}

func init() {
	rootCmd.AddCommand(delCmd)

	registerEditFlags(delCmd)
}
