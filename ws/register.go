package ws

import (
	. "codeit/model"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var languages = []Language{
	NewLanguage("c", ".c"),
	NewLanguage("cpp", ".cpp"),
	NewLanguage("golang", ".go"),
	NewLanguage("java", ".java"),
	NewLanguage("javascript", ".js"),
	NewLanguage("python2", ".py"),
	NewLanguage("python3", ".py"),
	NewLanguage("typescript", ".ts"),
}

func Register(apiGroup fiber.Router) {
	for _, l := range languages {
		apiGroup.Get(l.Name, websocket.New(l.Handler))
	}
}
