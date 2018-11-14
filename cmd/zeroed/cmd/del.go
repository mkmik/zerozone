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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("del called %q\n", viper.GetString("file"))
	},
}

func init() {
	rootCmd.AddCommand(delCmd)

	registerEditFlags(delCmd)
}
