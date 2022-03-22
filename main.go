package main

import (
	"codeit/api"
	"codeit/ws"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(cache.New())
	app.Static("/static/", "./static")

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./static/home.html")
	})

	app.Get("/editor", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./static/editor.html")
	})

	apiGroup := app.Group("/api")
	api.Register(apiGroup)

	wsGroup := app.Group("/ws")
	wsGroup.Use("/", func(ctx *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(ctx) {
			ctx.Locals("allowed", true)
			return ctx.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	ws.Register(wsGroup)

	log.Fatalln(app.Listen(":8080"))
}
