package web

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/gigapipehq/loggen/internal/cmd"
	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/web/utils"
)

func StartServer(ctx context.Context) error {
	app := fiber.New()
	app.Post("/", utils.ValidateRequest(&config.Config{}), func(ctx *fiber.Ctx) error {
		req := ctx.UserContext().Value("req").(*config.Config)
		if err := cmd.Do(req, "server request"); err != nil {
			return utils.Error(ctx, fiber.StatusInternalServerError, err.Error())
		}
		return utils.MessageResponse(ctx, fiber.StatusOK, "Done")
	})

	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Printf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	return app.Shutdown()
}
