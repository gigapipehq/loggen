package server

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/gigapipehq/loggen/web"
)

var serverCMD = &cobra.Command{
	Use:   "server",
	Short: "Run the generator in server-mode",
	Run: func(_ *cobra.Command, _ []string) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			exit := make(chan os.Signal, 1)
			signal.Notify(exit, os.Interrupt, os.Kill)
			<-exit
			cancel()
		}()

		if err := web.StartServer(ctx); err != nil {
			log.Printf("Error in webserver: %v", err)
		}
	},
}

func CMD() *cobra.Command {
	return serverCMD
}
