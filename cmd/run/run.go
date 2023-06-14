package run

import (
	"context"
	"fmt"
	"os"
	"time"

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
			if duration.Seconds() > 0 {
				cfg.Timeout = config.Duration(duration)
			}
			p := progress.NewBar(cfg.Rate*int(cfg.Timeout.Seconds()), os.Stdout)
			if err := cmd.Do(context.Background(), cfg, p); err != nil {
				fmt.Println(err)
			}
		},
	}
	duration time.Duration
)

func CMD() *cobra.Command {
	config.Load()
	cfg = config.Get()
	runCMD.Flags().StringVarP(
		&cfg.URL,
		"output",
		"o",
		cfg.URL,
		"output location for generated data. Either qryn URL or file location",
	)
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
		&duration,
		"timeout",
		"d",
		time.Duration(cfg.Timeout),
		"length of time to run the generator before exiting",
	)
	runCMD.Flags().StringVarP(
		&cfg.LogConfig.Format,
		"format",
		"f",
		cfg.LogConfig.Format,
		"format to use when sending logs",
	)
	runCMD.Flags().BoolVarP(
		&cfg.EnableMetrics,
		"enable-metrics",
		"m",
		cfg.EnableMetrics,
		"enable collection of Prometheus metrics",
	)
	return runCMD
}
