package controllers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/internal/generators"
)

func GetConfig(ctx *fiber.Ctx) error {
	args := ctx.Context().QueryArgs()
	catArgs := args.PeekMulti("category")
	categories := make([]string, len(catArgs))
	for i, category := range catArgs {
		categories[i] = string(category)
	}

	fromConfig := false
	if strings.ToLower(string(args.Peek("from_config"))) == "true" {
		fromConfig = true
	}
	return ctx.JSON(config.Get().LogConfig.Detailed(fromConfig, categories...))
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
