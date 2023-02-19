package controllers

import (
	"bufio"

	"github.com/gofiber/fiber/v2"

	"github.com/gigapipehq/loggen/internal/cmd"
	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/internal/progress"
	"github.com/gigapipehq/loggen/web/utils"
)

func Run(ctx *fiber.Ctx) error {
	req := ctx.UserContext().Value("req").(*config.Config)
	p := progress.NewServer(req.Rate * int(req.Timeout.Seconds()))

	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.Do(ctx.Context(), req, "server request", p)
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
}
