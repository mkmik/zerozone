package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	recordData string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new DNS record",
	Long:  `Add a new DNS record`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("adding %q to %q of type %q to file %q\n", recordData, recordName, recordData, viper.GetString("file"))
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	registerEditFlags(addCmd)

	addCmd.Flags().StringVarP(&recordData, "data", "d", "", "record data")
	addCmd.MarkFlagRequired("data")

	viper.BindPFlag("file", addCmd.Flags().Lookup("file"))
}
