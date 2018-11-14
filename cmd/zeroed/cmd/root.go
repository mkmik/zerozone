package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	apiAddrCfg        = "api"
	pubKeyCfg         = "pubkey"
	zeroZoneDomainCfg = "zeroZoneDomain"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "zeroed",
	Short: "Zero Zone editor",
	Long: `zeroed allows you to manage your private "zero zone", by manipulating a JSON zone
file and publishing it via IPFS`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.zeroed.yaml)")

	rootCmd.PersistentFlags().String(apiAddrCfg, "/ip4/127.0.0.1/tcp/5001", "address of ipfs api server")
	rootCmd.PersistentFlags().String(pubKeyCfg, "self", "pubkey name of hash")
	rootCmd.PersistentFlags().String(zeroZoneDomainCfg, "0zone.mkm.pub", "domain name of the ZeroZone service")

	viper.BindPFlags(rootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		u, err := user.Current()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".zeroed" (without extension).
		viper.AddConfigPath(u.HomeDir)
		viper.SetConfigName(".zeroed")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
