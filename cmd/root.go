package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "t2q2t",
	Short: "tcp/quic port forward tool",
	Long: `tcp/quic port forward tool
  t2q2t <subcommand> <Listen Addr> <forward Addr>  

  go run ./t2q2t.go t2q 0.0.0.0:2022 127.0.0.1:2022
  go run ./t2q2t.go q2t 0.0.0.0:2022 127.0.0.1:22
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".t2q2t")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
