package api

import "github.com/gofiber/fiber/v2"

func Register(apiGroup fiber.Router) {
	apiGroup.Post("/c", HandleC)
	apiGroup.Post("/cpp", HandleCpp)
	apiGroup.Post("/java", HandleJava)
	apiGroup.Post("/go", HandleGo)
	apiGroup.Post("/python2", HandlePython2)
	apiGroup.Post("/python3", HandlePython3)
}
