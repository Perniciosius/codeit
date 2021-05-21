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

func HandlePython3(ctx *fiber.Ctx) error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return ctx.JSON(map[string]string{
			"output": "Some error has occurred. Please try again after sometime",
		})
	}

	folderName := fmt.Sprintf("python3_%v", time.Now().UnixNano())
	err = os.MkdirAll(folderName, 0775)
	if err != nil {
		log.Println(err)
		return ctx.JSON(map[string]string{
			"output": "Some error has occurred. Please try again after sometime",
		})
	}

	defer utilities.Cleanup(folderName)

	executeCommand := "python3 main.py"
	dockerCommand := fmt.Sprintf("docker run --rm -v %v/%v:/work -w /work python bash script.sh", pwd, folderName)

	// parse request body
	compileRequestBody := new(model.CompileRequestBody)
	_ = ctx.BodyParser(compileRequestBody)

	script := utilities.BuildScript("", executeCommand, 2)

	fileName := fmt.Sprintf("%v/main.py", folderName)
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
