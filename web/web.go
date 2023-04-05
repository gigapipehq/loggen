package web

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mailru/easyjson"

	"github.com/gigapipehq/loggen/web/config"
	"github.com/gigapipehq/loggen/web/run"
)

var routers = map[string]func(router fiber.Router){
	"/run":    run.Register,
	"/config": config.Register,
}

func StartServer(ctx context.Context) error {
	app := fiber.New(fiber.Config{
		JSONEncoder: func(v interface{}) ([]byte, error) {
			return easyjson.Marshal(v.(easyjson.Marshaler))
		},
		JSONDecoder: func(data []byte, v interface{}) error {
			return easyjson.Unmarshal(data, v.(easyjson.Unmarshaler))
		},
	})

	app.Get("/status", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	})

	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Context().SetContentType("application/json")
		return ctx.Next()
	})

	for path, register := range routers {
		group := app.Group(path)
		register(group)
	}

	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Printf("server error: %v", err)
		}
	}()
	<-ctx.Done()
	return app.Shutdown()
}
