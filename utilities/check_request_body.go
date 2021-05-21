package utilities

import (
	"codeit/model"
	"github.com/gofiber/fiber/v2"
)

func CheckRequestBody(ctx *fiber.Ctx) error {
	body := new(model.CompileRequestBody)
	if err := ctx.BodyParser(body); err != nil || body.Code == "" {
		return ctx.JSON(map[string]string{
			"output": "Empty program",
		})
	}
	return ctx.Next()
}
