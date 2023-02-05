package web

import (
	"bufio"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/gigapipehq/loggen/internal/cmd"
	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/internal/progress"
	"github.com/gigapipehq/loggen/web/utils"
)

func StartServer(ctx context.Context) error {
	app := fiber.New()
	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Context().SetContentType("application/json")
		return ctx.Next()
	})
	app.Post("/", utils.ValidateRequest(&config.Config{}), func(ctx *fiber.Ctx) error {
		req := ctx.UserContext().Value("req").(*config.Config)
		p := progress.NewServer(req.Rate * int(req.Timeout.Seconds()))

		errCh := make(chan error, 1)
		go func() {
			errCh <- cmd.Do(req, "server request", p)
		}()
		ctx.Response().Header.Set("Cache-Control", "no-cache")
		ctx.Response().Header.Set("Connection", "keep-alive")
		ctx.Response().Header.Set("Transfer-Encoding", "chunked")
		ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			defer close(errCh)
			p.WriteProgress(w)
			if err := <-errCh; err != nil {
				_, _ = w.Write(utils.ErrorResponseBytes(err))
				_ = w.Flush()
			}
		})
		return nil
	})

	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Printf("server error: %v", err)
		}
	}()
	<-ctx.Done()
	return app.Shutdown()
}
