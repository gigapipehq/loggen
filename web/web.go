package web

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mailru/easyjson"

	"github.com/gigapipehq/loggen/internal/cmd"
	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/internal/progress"
	"github.com/gigapipehq/loggen/web/utils"
)

func StartServer(ctx context.Context) error {
	app := fiber.New(fiber.Config{
		JSONEncoder: func(v interface{}) ([]byte, error) {
			return easyjson.Marshal(v.(easyjson.Marshaler))
		},
		JSONDecoder: func(data []byte, v interface{}) error {
			return easyjson.Unmarshal(data, v.(easyjson.Unmarshaler))
		},
	})
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
	app.Get("/config", func(ctx *fiber.Ctx) error {
		return ctx.JSON(config.Get())
	})
	app.Get("/config/logs", func(ctx *fiber.Ctx) error {
		q := ctx.Context().QueryArgs().PeekMulti("category")
		categories := make([]string, len(q))
		for i, category := range q {
			categories[i] = string(category)
		}
		return ctx.JSON(config.Get().LogConfig.Detailed(categories...))
	})
	app.Patch("/config", func(ctx *fiber.Ctx) error {
		var cfg map[string]interface{}
		if err := json.Unmarshal(ctx.Body(), &cfg); err != nil {
			return ctx.Status(fiber.StatusBadRequest).Send(utils.ErrorResponseBytes(err))
		}
		for k, v := range cfg {
			if err := config.UpdateSettingValue(k, fmt.Sprintf("%v", v)); err != nil {
				return ctx.Status(fiber.StatusBadRequest).Send(utils.ErrorResponseBytes(err))
			}
		}
		return ctx.SendStatus(fiber.StatusNoContent)
	})

	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Printf("server error: %v", err)
		}
	}()
	<-ctx.Done()
	return app.Shutdown()
}
