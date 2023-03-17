package cmd

import (
	"github.com/spf13/cobra"

	configcmd "github.com/gigapipehq/loggen/cmd/config"
	"github.com/gigapipehq/loggen/cmd/lambda"
	"github.com/gigapipehq/loggen/cmd/run"
	"github.com/gigapipehq/loggen/cmd/server"
	"github.com/gigapipehq/loggen/internal/config"
)

var rootCMD = &cobra.Command{
	Use:   "loggen",
	Short: "A fake log, metric and trace generator for qryn Cloud",
}

func init() {
	config.Load()
	rootCMD.AddCommand(configcmd.CMD(), lambda.CMD(), run.CMD(), server.CMD())
}

func Execute() error {
	return rootCMD.Execute()
}
