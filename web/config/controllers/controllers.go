package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/web/utils"
)

func GetConfig(ctx *fiber.Ctx) error {
	return ctx.JSON(config.Get())
}

func UpdateConfig(ctx *fiber.Ctx) error {
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
}
