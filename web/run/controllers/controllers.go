package controllers

import (
	"bufio"
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/gigapipehq/loggen/internal/cmd"
	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/internal/progress"
	"github.com/gigapipehq/loggen/web/utils"
)

func Run(ctx *fiber.Ctx, cfg *config.Config) error {
	p := progress.NewServer(cfg.Rate * int(cfg.Timeout.Seconds()))

	errCh := make(chan error, 1)
	cmdCtx, cancel := context.WithCancel(ctx.Context())
	go func() {
		errCh <- cmd.Do(cmdCtx, cfg, p)
	}()
	ctx.Response().Header.Set("Cache-Control", "no-cache")
	ctx.Response().Header.Set("Connection", "keep-alive")
	ctx.Response().Header.Set("Transfer-Encoding", "chunked")
	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer close(errCh)
		defer cancel()

		p.WriteProgress(w)
		if err := <-errCh; err != nil {
			_, _ = w.Write(utils.ErrorResponseBytes(err))
			_ = w.Flush()
		}
	})
	return nil
}
