package api

import (
	"codeit/model"
	"codeit/utilities"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

func HandleJava(ctx *fiber.Ctx) error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	folderName := fmt.Sprintf("java_%v", time.Now().UnixNano())
	err = os.MkdirAll(folderName, 0775)
	if err != nil {
		log.Fatalln(err)
	}

	compileCommand := "javac main.java"
	executeCommand := "java main"
	dockerCommand := fmt.Sprintf("docker run -t --rm -v %v/%v:/work -w /work openjdk sh script.sh", pwd, folderName)

	// parse request body
	compileRequestBody := new(model.CompileRequestBody)
	_ = ctx.BodyParser(compileRequestBody)

	if compileRequestBody.Arguments != "" {
		compileCommand = fmt.Sprintf("%v %v", compileCommand, compileRequestBody.Arguments)
	}

	script := []byte(fmt.Sprintf("#!/bin/bash\n%v && %v", compileCommand, executeCommand))

	fileName := fmt.Sprintf("%v/main.java", folderName)
	err = ioutil.WriteFile(fileName, []byte(compileRequestBody.Code), 0664)
	if err != nil {
		log.Fatalln(err)
	}
	scriptName := fmt.Sprintf("%v/script.sh", folderName)
	err = ioutil.WriteFile(scriptName, script, 0664)
	if err != nil {
		log.Fatalln(err)
	}

	cmd := exec.Command("sh", "-c", dockerCommand)

	var output []byte
	output, err = cmd.CombinedOutput()
	if err != nil && len(output) < 1 {
		log.Fatalln(err)
	}

	defer utilities.Cleanup(folderName)

	return ctx.JSON(map[string]string{
		"output": string(output),
	})
}
