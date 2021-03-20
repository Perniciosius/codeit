package api

import (
	"codeit/model"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func HandlePython2(ctx *fiber.Ctx) error {
	executeCommand := "python2 main.py"
	compileRequestBody := new(model.CompileRequestBody)
	err := ctx.BodyParser(compileRequestBody)
	if err != nil {
		log.Fatalln(err)
	}
	folderName := fmt.Sprintf("python2_%v", time.Now().UnixNano())
	//err = os.MkdirAll(folderName, 1775)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	if compileRequestBody.Code == "" {
		return ctx.JSON(map[string]string{
			"error": "Empty program",
		})
	}

	if compileRequestBody.Arguments != "" {
		executeCommand = fmt.Sprintf("%v %v", executeCommand, compileRequestBody.Arguments)
	}

	return ctx.JSON(map[string]string{
		"folder_name":     folderName,
		"execute_command": executeCommand,
		"code":            compileRequestBody.Code,
	})
}
