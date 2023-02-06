package reset

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gigapipehq/loggen/internal/config"
)

var resetCMD = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to initial defaults",
	Run: func(_ *cobra.Command, args []string) {
		if err := config.Reset(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Configuration restored to initial defaults updated")
	},
}

func CMD() *cobra.Command {
	return resetCMD
}
