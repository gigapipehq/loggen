package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/internal/generators"
)

func GetConfig(ctx *fiber.Ctx) error {
	q := ctx.Context().QueryArgs().PeekMulti("category")
	categories := make([]string, len(q))
	for i, category := range q {
		categories[i] = string(category)
	}
	return ctx.JSON(config.Get().LogConfig.Detailed(categories...))
}

func GetExampleLine(ctx *fiber.Ctx) error {
	lc := ctx.UserContext().Value("req").(*config.LogConfig)
	example := generators.Example(*lc)
	if lc.Format == "logfmt" {
		example = []byte(fmt.Sprintf("\"%s\"", example))
	}
	resp := []byte(fmt.Sprintf("{\"line\":%s}", example))
	return ctx.Send(resp)
}
