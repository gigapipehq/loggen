package get

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mailru/easyjson"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/gigapipehq/loggen/internal/config"
)

type outputFmt string

const (
	outputFmtYAML outputFmt = "yaml"
	outputFmtJSON outputFmt = "json"
)

func (e *outputFmt) String() string {
	return string(*e)
}

func (e *outputFmt) Set(v string) error {
	for _, f := range validOutputFormats {
		if f == v {
			*e = outputFmt(v)
			return nil
		}
	}
	prefix := "must be"
	if len(validOutputFormats) == 2 {
		return fmt.Errorf("%s either %s", prefix, outputFormatHelpList())
	}
	return fmt.Errorf("%s %s", prefix, outputFormatHelpList())
}

func (e *outputFmt) Type() string {
	return "string"
}

var (
	getCMD = &cobra.Command{
		Use:       "get [setting name]",
		Short:     "Show current default configuration settings",
		ValidArgs: config.SettingNames(),
		Args:      cobra.MatchAll(cobra.MaximumNArgs(1), config.ValidArgSettingName("get")),
		Run: func(_ *cobra.Command, args []string) {
			cfg := config.Get()
			if len(args) == 0 {
				switch outputFormat {
				case outputFmtYAML:
					b, _ := yaml.Marshal(cfg)
					fmt.Print(string(b))
				case outputFmtJSON:
					var out bytes.Buffer
					b, _ := easyjson.Marshal(cfg)
					_ = json.Indent(&out, b, "", "\t")
					fmt.Println(out.String())
				}
				return
			}
			fmt.Printf("%v\n", config.GetSettingValue(args[0]))
		},
	}
	validOutputFormats = []string{"yaml", "json"}
	outputFormat       = outputFmtYAML
)

func CMD() *cobra.Command {
	getCMD.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ToLower(name))
	})
	ofusage := fmt.Sprintf("output format of config; %s", outputFormatHelpList())
	getCMD.Flags().VarP(&outputFormat, "output-format", "f", ofusage)
	_ = getCMD.RegisterFlagCompletionFunc(
		"output-format",
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return validOutputFormats, cobra.ShellCompDirectiveNoFileComp
		},
	)

	return getCMD
}

func outputFormatHelpList() string {
	if len(validOutputFormats) == 2 {
		return fmt.Sprintf("%s, or %s", validOutputFormats[0], validOutputFormats[1])
	}
	return fmt.Sprintf(
		"one of %s, or %s",
		strings.Join(validOutputFormats, ", "), validOutputFormats[len(validOutputFormats)-1],
	)
}
