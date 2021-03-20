package api

import (
	"codeit/model"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func HandleJava(ctx *fiber.Ctx) error {
	compileCommand := "javac main.java"
	executeCommand := "java main"
	compileRequestBody := new(model.CompileRequestBody)
	err := ctx.BodyParser(compileRequestBody)
	if err != nil {
		log.Fatalln(err)
	}
	folderName := fmt.Sprintf("java_%v", time.Now().UnixNano())
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
		"compile_command": compileCommand,
		"execute_command": executeCommand,
		"code":            compileRequestBody.Code,
	})
}
