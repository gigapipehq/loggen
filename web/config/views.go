package config

import (
	"github.com/gofiber/fiber/v2"

	"github.com/gigapipehq/loggen/web/config/controllers"
	"github.com/gigapipehq/loggen/web/config/logs"
)

var routers = map[string]func(router fiber.Router){
	"/logs": logs.Register,
}

func Register(router fiber.Router) {
	router.Get("/", controllers.GetConfig)
	router.Patch("/", controllers.UpdateConfig)

	for path, register := range routers {
		group := router.Group(path)
		register(group)
	}
}
