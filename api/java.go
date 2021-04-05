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
	"regexp"
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

	defer utilities.Cleanup(folderName)

	dockerCommand := fmt.Sprintf("docker run --rm -v %v/%v:/work -w /work openjdk bash script.sh", pwd, folderName)

	// parse request body
	compileRequestBody := new(model.CompileRequestBody)
	_ = ctx.BodyParser(compileRequestBody)

	re := regexp.MustCompile(`.*class\s+([A-Za-z0-9_]+)`)
	var className string
	if match := re.FindStringSubmatch(compileRequestBody.Code); len(match) < 2 {
		className = "main"
	} else {
		className = match[1]
	}

	compileCommand := fmt.Sprintf("javac %v.java", className)
	executeCommand := fmt.Sprintf("java %v", className)

	if compileRequestBody.CompileArguments != "" {
		compileCommand = fmt.Sprintf("%v %v", compileCommand, compileRequestBody.CompileArguments)
	}

	if compileRequestBody.RuntimeArguments != "" {
		executeCommand = fmt.Sprintf("%v %v", executeCommand, compileRequestBody.RuntimeArguments)
	}

	script := utilities.BuildScript(compileCommand, executeCommand, 1)

	fileName := fmt.Sprintf("%v/%v.java", folderName, className)
	err = ioutil.WriteFile(fileName, []byte(compileRequestBody.Code), 0664)
	if err != nil {
		log.Fatalln(err)
	}

	inputFileName := fmt.Sprintf("%v/input.txt", folderName)
	err = ioutil.WriteFile(inputFileName, []byte(compileRequestBody.Input), 0664)
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

	return ctx.JSON(map[string]string{
		"output": string(output),
	})
}
