package set

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gigapipehq/loggen/internal/config"
)

var setCMD = &cobra.Command{
	Use:       "set [setting name] [setting value]",
	Short:     "Set current default configuration setting",
	ValidArgs: config.SettingNames(),
	Args:      cobra.MatchAll(cobra.ExactArgs(2), config.ValidArgSettingName("set")),
	Run: func(_ *cobra.Command, args []string) {
		if err := config.UpdateSettingValue(args[0], args[1]); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Setting %s updated\n", args[0])
	},
}

func CMD() *cobra.Command {
	return setCMD
}
