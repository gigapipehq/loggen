package utils

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type UnprocessableEntityResponse struct {
	Errors map[string]unprocessableEntityError `json:"errors"`
}

type unprocessableEntityError struct {
	Validator string `json:"validator"`
	Value     string `json:"value"`
}

var requestValidator = validator.New()

func MessageResponse(ctx *fiber.Ctx, status int, message string) error {
	return ctx.Status(status).JSON(map[string]string{
		"message": message,
	})
}

func Error(ctx *fiber.Ctx, status int, error string) error {
	return ctx.Status(status).JSON(ErrorResponse{
		Error: error,
	})
}

func ValidateRequest(req any) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return validateRequest(ctx, req)
	}
}

func validateRequest(ctx *fiber.Ctx, req any) error {
	if string(ctx.Request().Header.ContentType()) != "application/json" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid Content-Type. Use application/json",
		})
	}
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	if err := requestValidator.Struct(req); err != nil {
		errors := map[string]unprocessableEntityError{}
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Namespace()] = unprocessableEntityError{
				Validator: err.ActualTag(),
				Value:     err.Param(),
			}
			err.Type()
		}
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(UnprocessableEntityResponse{
			Errors: errors,
		})
	}
	uctx := context.WithValue(ctx.UserContext(), "req", req)
	ctx.SetUserContext(uctx)
	return ctx.Next()
}
