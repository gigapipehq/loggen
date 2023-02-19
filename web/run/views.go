package run

import (
	"github.com/gofiber/fiber/v2"

	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/web/run/controllers"
	"github.com/gigapipehq/loggen/web/utils"
)

func Register(router fiber.Router) {
	router.Post("/", utils.ValidateRequest(&config.Config{}), controllers.Run)
}
