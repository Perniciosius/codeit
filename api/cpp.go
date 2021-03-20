package api

import (
	"codeit/model"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func HandleCpp(ctx *fiber.Ctx) error {
	compileCommand := "g++ main.c -o main"
	executeCommand := "./main"
	compileRequestBody := new(model.CompileRequestBody)
	err := ctx.BodyParser(compileRequestBody)
	if err != nil {
		log.Fatalln(err)
	}
	folderName := fmt.Sprintf("cpp_%v", time.Now().UnixNano())
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
		compileCommand = fmt.Sprintf("%v %v", compileCommand, compileRequestBody.Arguments)
	}

	return ctx.JSON(map[string]string{
		"file_name":       folderName,
		"compile_command": compileCommand,
		"execute_command": executeCommand,
		"code":            compileRequestBody.Code,
	})
}
