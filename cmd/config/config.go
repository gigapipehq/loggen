package config

import (
	"github.com/spf13/cobra"

	"github.com/gigapipehq/loggen/cmd/config/get"
	"github.com/gigapipehq/loggen/cmd/config/set"
)

var configCMD = &cobra.Command{
	Use:   "config",
	Short: "Configuration commands",
}

func init() {
	configCMD.AddCommand(get.CMD(), set.CMD())
}

func CMD() *cobra.Command {
	return configCMD
}
