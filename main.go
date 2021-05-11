package main

import (
	"codeit/api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(cache.New())
	app.Static("/static", "./static")

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("static/home.html")
	})

	app.Get("/editor", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("static/editor.html")
	})

	apiGroup := app.Group("/api")
	api.Register(apiGroup)

	log.Fatalln(app.Listen(":8080"))
}
