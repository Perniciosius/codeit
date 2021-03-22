package api

import (
	"codeit/utilities"
	"github.com/gofiber/fiber/v2"
)

func Register(apiGroup fiber.Router) {
	apiGroup.Post("/c", utilities.CheckRequestBody, HandleC)
	apiGroup.Post("/cpp", utilities.CheckRequestBody, HandleCpp)
	apiGroup.Post("/java", utilities.CheckRequestBody, HandleJava)
	apiGroup.Post("/go", utilities.CheckRequestBody, HandleGo)
	apiGroup.Post("/python2", utilities.CheckRequestBody, HandlePython2)
	apiGroup.Post("/python3", utilities.CheckRequestBody, HandlePython3)

	apiGroup.Use(func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(404)
	})
}
