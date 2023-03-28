package controllers

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"

	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/internal/generators"
	"github.com/gigapipehq/loggen/web/config/logs/schemas"
)

func GetConfig(ctx *fiber.Ctx) error {
	args := ctx.Context().QueryArgs()
	catArgs := args.PeekMulti("categories")
	categories := make(map[string]struct{}, len(catArgs))
	for _, category := range catArgs {
		categories[string(category)] = struct{}{}
	}

	var funcs schemas.GeneratorFunctionList
	for _, info := range gofakeit.FuncLookups {
		if _, ok := categories[info.Category]; ok {
			funcs = append(funcs, schemas.GeneratorFunction{
				Display:     info.Display,
				Category:    info.Category,
				Description: info.Description,
				Example:     info.Example,
				Params:      info.Params,
			})
		}
	}
	return ctx.JSON(funcs)
}

func GetCategories(ctx *fiber.Ctx) error {
	categories := map[string]struct{}{}
	for _, info := range gofakeit.FuncLookups {
		if _, ok := categories[info.Category]; !ok {
			categories[info.Category] = struct{}{}
		}
	}
	s := make(schemas.CategoryList, 0, len(categories))
	for c := range categories {
		s = append(s, c)
	}
	return ctx.JSON(s)
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
