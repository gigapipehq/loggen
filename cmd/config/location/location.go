package location

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gigapipehq/loggen/internal/config"
)

var (
	locationCMD = &cobra.Command{
		Use:   "location",
		Short: "print out configuration file location",
		Run: func(_ *cobra.Command, args []string) {
			p := config.ConfigFilename
			if dirOnly {
				p = config.BasePath
			}
			fmt.Println(p)
		},
	}
	dirOnly bool
)

func CMD() *cobra.Command {
	locationCMD.Flags().BoolVar(
		&dirOnly, "dir-only", false, "only print out configuration file directory",
	)
	return locationCMD
}
