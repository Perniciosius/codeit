package api

import (
	"codeit/utilities"
	"github.com/gofiber/fiber/v2"
)

func Register(apiGroup fiber.Router) {
	apiGroup.Post("/c", utilities.CheckRequestBody, HandleC)
	apiGroup.Post("/cpp", utilities.CheckRequestBody, HandleCpp)
	apiGroup.Post("/go", utilities.CheckRequestBody, HandleGo)
	apiGroup.Post("/java", utilities.CheckRequestBody, HandleJava)
	apiGroup.Post("/javascript", utilities.CheckRequestBody, HandleJavascript)
	apiGroup.Post("/python", utilities.CheckRequestBody, HandlePython3)
	apiGroup.Post("/python2", utilities.CheckRequestBody, HandlePython2)
	apiGroup.Post("/python3", utilities.CheckRequestBody, HandlePython3)
	apiGroup.Post("typescript", utilities.CheckRequestBody, HandleTypescript)

	apiGroup.Use(func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(404)
	})
}
