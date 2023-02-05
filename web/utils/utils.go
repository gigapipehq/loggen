package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type unprocessableEntityResponse struct {
	Errors map[string]unprocessableEntityError `json:"errors"`
}

type unprocessableEntityError struct {
	Validator string `json:"validator"`
	Value     string `json:"value"`
}

func SendError(ctx *fiber.Ctx, err error) error {
	if _, werr := ctx.Write(ErrorResponseBytes(err)); werr != nil {
		return werr
	}
	return nil
}

func ErrorResponseBytes(err error) []byte {
	return []byte(fmt.Sprintf("{\"error\": \"%s\"}", err.Error()))
}

var requestValidator = validator.New()

func ValidateRequest(req any) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return validateRequest(ctx, req)
	}
}

func validateRequest(ctx *fiber.Ctx, req any) error {
	if string(ctx.Request().Header.ContentType()) != "application/json" {
		return SendError(ctx, errors.New("Invalid Content-Type. Use application/json"))
	}
	if err := ctx.BodyParser(req); err != nil {
		return SendError(ctx, err)
	}

	if err := requestValidator.Struct(req); err != nil {
		uees := map[string]unprocessableEntityError{}
		for _, err := range err.(validator.ValidationErrors) {
			uees[err.Namespace()] = unprocessableEntityError{
				Validator: err.ActualTag(),
				Value:     err.Param(),
			}
			err.Type()
		}
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(unprocessableEntityResponse{
			Errors: uees,
		})
	}
	uctx := context.WithValue(ctx.UserContext(), "req", req)
	ctx.SetUserContext(uctx)
	return ctx.Next()
}
