package logs

import (
	"github.com/gofiber/fiber/v2"

	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/web/config/logs/controllers"
	"github.com/gigapipehq/loggen/web/utils"
)

func Register(router fiber.Router) {
	router.Get("/", controllers.GetConfig)
	router.Get("/categories", controllers.GetCategories)
	router.Post("/example", utils.ValidateRequest(&config.LogConfig{}), controllers.GetExampleLine)
}
