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

func HandleCpp(ctx *fiber.Ctx) error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return ctx.JSON(map[string]string{
			"output": "Some error has occurred. Please try again after sometime",
		})
	}

	folderName := fmt.Sprintf("cpp_%v", time.Now().UnixNano())
	err = os.MkdirAll(folderName, 0775)
	if err != nil {
		log.Println(err)
		return ctx.JSON(map[string]string{
			"output": "Some error has occurred. Please try again after sometime",
		})
	}

	defer utilities.Cleanup(folderName)

	compileCommand := "g++ main.cpp -o main"
	executeCommand := "./main"
	dockerCommand := fmt.Sprintf("docker run --rm -v %v/%v:/work -w /work gcc bash script.sh", pwd, folderName)

	// parse request body
	compileRequestBody := new(model.CompileRequestBody)
	_ = ctx.BodyParser(compileRequestBody)

	if compileRequestBody.CompileArguments != "" {
		compileCommand = fmt.Sprintf("%v %v", compileCommand, compileRequestBody.CompileArguments)
	}

	if compileRequestBody.RuntimeArguments != "" {
		executeCommand = fmt.Sprintf("%v %v", executeCommand, compileRequestBody.RuntimeArguments)
	}

	script := utilities.BuildScript(compileCommand, executeCommand, 1)

	fileName := fmt.Sprintf("%v/main.cpp", folderName)
	err = ioutil.WriteFile(fileName, []byte(compileRequestBody.Code), 0664)
	if err != nil {
		log.Println(err)
		return ctx.JSON(map[string]string{
			"output": "Some error has occurred. Please try again after sometime",
		})
	}

	inputFileName := fmt.Sprintf("%v/input.txt", folderName)
	err = ioutil.WriteFile(inputFileName, []byte(compileRequestBody.Input), 0664)
	if err != nil {
		log.Println(err)
		return ctx.JSON(map[string]string{
			"output": "Some error has occurred. Please try again after sometime",
		})
	}

	scriptName := fmt.Sprintf("%v/script.sh", folderName)
	err = ioutil.WriteFile(scriptName, script, 0664)
	if err != nil {
		log.Println(err)
		return ctx.JSON(map[string]string{
			"output": "Some error has occurred. Please try again after sometime",
		})
	}

	cmd := exec.Command("sh", "-c", dockerCommand)

	var output []byte
	output, err = cmd.CombinedOutput()
	if err != nil && len(output) < 1 {
		log.Println(err)
		return ctx.JSON(map[string]string{
			"output": "Some error has occurred. Please try again after sometime",
		})
	}

	return ctx.JSON(map[string]string{
		"output": string(output),
	})
}
