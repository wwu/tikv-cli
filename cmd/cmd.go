package cmd

import (
	"fmt"
	"os"

	pingcaplog "github.com/pingcap/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	// Host is the PD host address.
	Host string
	// Port is the PD port.
	Port string
	// Mode is the client mode: raw/txn
	Mode string
	// Keyspace name for tikv V2 storage
	KeySpace string
	// APIVersion is the API version: v1/v1ttl/v2
	APIVersion string
	// Debug determines whether to enable logging in tikv/client-go.
	Debug bool
)

var rootCmd = &cobra.Command{
	Use:   "tikv-cli",
	Long:  `A CLI for TiKV cluster through PD. You can enter the interactive shell by root command.`,
	Short: "Interact with TiKV cluster through PD",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) (err error) {
		if cmd.Name() == "help" || cmd.Name() == "version" {
			return
		}

		if !Debug {
			// Disable logging in tikv/client-go
			pingcaplog.ReplaceGlobals(zap.NewNop(), nil)
		}

		c, err = newClient()
		return
	},
	RunE: shellRunE,
	PersistentPostRunE: func(cmd *cobra.Command, _ []string) error {
		if c != nil {
			return c.Close(cmd.Context())
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&Host, "host", "h", "localhost", "PD host address")
	rootCmd.PersistentFlags().StringVarP(&Port, "port", "p", "2379", "PD port")
	rootCmd.PersistentFlags().StringVarP(&Mode, "mode", "m", "txn", "Client mode. raw/txn")
	rootCmd.PersistentFlags().StringVarP(&KeySpace, "keyspace", "k", "", "Tikv keyspace, default is empty")
	rootCmd.PersistentFlags().StringVarP(&APIVersion, "api-version", "a", "v2", "API version. v1/v1ttl/v2")
	rootCmd.PersistentFlags().Bool("help", false, "Help for tikv-cli")
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Debug determines whether to enable logging in tikv/client-go")

	rootCmd.AddCommand(versionCmd, putCmd, getCmd, deleteCmd, ttlCmd, scanCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
