package config

import (
	"github.com/spf13/cobra"

	"github.com/gigapipehq/loggen/cmd/config/get"
	"github.com/gigapipehq/loggen/cmd/config/location"
	"github.com/gigapipehq/loggen/cmd/config/reset"
	"github.com/gigapipehq/loggen/cmd/config/set"
	"github.com/gigapipehq/loggen/internal/config"
)

var configCMD = &cobra.Command{
	Use:   "config",
	Short: "Configuration commands",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		config.Load()
	},
}

func init() {
	configCMD.AddCommand(get.CMD(), location.CMD(), set.CMD(), reset.CMD())
}

func CMD() *cobra.Command {
	return configCMD
}
