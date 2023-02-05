package run

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/gigapipehq/loggen/internal/cmd"
	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/internal/progress"
)

var (
	cfg    *config.Config
	runCMD = &cobra.Command{
		Use:   "run",
		Short: "Run the generator in cli-mode",
		Run: func(_ *cobra.Command, _ []string) {
			p := progress.NewBar(cfg.Rate*int(cfg.Timeout.Seconds()), os.Stdout)
			if err := cmd.Do(cfg, "run in cli-mode", p); err != nil {
				fmt.Println(err)
			}
		},
	}
)

func CMD() *cobra.Command {
	cfg = config.Get()
	runCMD.Flags().StringVarP(
		&cfg.APIKey,
		"api-key",
		"k",
		cfg.APIKey,
		"API key to use for authenticating with qryn Cloud",
	)
	runCMD.Flags().StringVarP(
		&cfg.APISecret,
		"api-secret",
		"s",
		cfg.APISecret,
		"API key to use for authenticating with qryn Cloud",
	)
	runCMD.Flags().StringToStringVarP(
		&cfg.Labels,
		"labels",
		"l",
		cfg.Labels,
		"labels for each log",
	)
	runCMD.Flags().IntVarP(
		&cfg.Rate,
		"rate",
		"r",
		cfg.Rate,
		"number of logs to generate per second",
	)
	runCMD.Flags().DurationVarP(
		&cfg.Timeout,
		"timeout",
		"t",
		cfg.Timeout,
		"length of time to run the generator before exiting",
	)
	return runCMD
}
