package model

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/gofiber/websocket/v2"
)

type Language struct {
	Name      string
	Filename  string
	Extension string
}

func NewLanguage(name, extension string) Language {
	return Language{
		Name:      name,
		Filename:  "main",
		Extension: extension,
	}
}

func (l *Language) getCommand(body *WsBody) []string {
	if l.Name == "java" {
		pattern := regexp.MustCompile(`(?s)class\s+(\w+).*?public\s+static\s+void\s+main\s*\(String(?:\s*\[\]\s+\w+|\s+\w+\s*\[\])\)`)
		matches := pattern.FindStringSubmatch(body.Code)
		if len(matches) >= 2 {
			l.Filename = matches[1]
		}
	}

	file := fmt.Sprintf("%v.%v", l.Filename, l.Extension)
	command := []string{}
	switch l.Name {
	case "c":
		command = append(command, "gcc", file, "-o", "main", "&&", "./main")
	case "cpp":
		command = append(command, "g++", file, "-o", "main", "&&", "./main")
	case "golang":
		command = append(command, "go", "run", file)
	case "java":
		command = append(command, "javac", file, "&&", "java", l.Filename)
	case "javascript":
		command = append(command, "node", file)
	case "python2":
		command = append(command, "python2", file)
	case "python3":
		command = append(command, "python3", file)
	case "typescript":
		command = append(command, "tsc", file, "&&", "node", l.Filename+".js")
	}

	return command
}

func (l *Language) Handler(c *websocket.Conn) {

	var (
		msg []byte
		err error
	)

	errorResponse := []byte("Server Error")

	if _, msg, err = c.ReadMessage(); err != nil {
		log.Println(err)

	}

	body := new(WsBody)
	if err := json.Unmarshal(msg, body); err != nil {
		log.Println(err)
	}

	directory, err := os.MkdirTemp("", "")
	if err != nil {
		log.Println(err)
		c.WriteMessage(1, errorResponse)
		return
	}
	defer os.RemoveAll(directory)

	if err := os.WriteFile(fmt.Sprintf("%v/%v.%v", directory, l.Filename, l.Extension), []byte(body.Code), 0664); err != nil {
		log.Println(err)
		c.WriteMessage(1, errorResponse)
		return
	}

	var dockerClient *client.Client
	dockerClient, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Println(err)
		c.WriteMessage(1, errorResponse)
		return
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var resp container.ContainerCreateCreatedBody
	resp, err = dockerClient.ContainerCreate(timeoutContext, &container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		Tty:          true,
		Image:        l.Name,
		WorkingDir:   "/work",
		Cmd:          l.getCommand(body),
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: directory,
				Target: "/work",
			},
		},
	}, nil, nil, "")

	if err != nil {
		log.Println(err)
		c.WriteMessage(1, errorResponse)
		return
	}

	if err = dockerClient.ContainerStart(timeoutContext, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Println(err)
		c.WriteMessage(1, errorResponse)
		return
	}

	var containerResp types.HijackedResponse
	containerResp, err = dockerClient.ContainerAttach(timeoutContext, resp.ID, types.ContainerAttachOptions{
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Stream: true,
	})
	if err != nil {
		log.Println(err)
		c.WriteMessage(1, errorResponse)
		return
	}
	defer containerResp.Close()

	bufin := bufio.NewReader(containerResp.Reader)
	inout := make(chan []byte)
	output := make(chan []byte)

	// Write to docker container
	go func(w io.WriteCloser) {
		for {
			data, ok := <-inout
			if !ok {
				w.Close()
				return
			}
			w.Write(append(data, '\n'))
		}
	}(containerResp.Conn)

	// Receive from docker container
	go func() {
		for {
			buffer := make([]byte, 4096, 4096)
			c, err := bufin.Read(buffer)
			if err != nil {
				fmt.Println(err)
			}
			if c > 0 {
				output <- buffer[:c]
			}
			if c == 0 {
				output <- []byte{' '}
			}
			if err != nil {
				break
			}
		}
	}()

	// Connect STDOUT to websocket
	go func() {
		data := <-output
		stringData := string(data[:])
		if !utf8.ValidString(stringData) {
			v := make([]rune, 0, len(stringData))
			for i, r := range stringData {
				if r == utf8.RuneError {
					_, size := utf8.DecodeRuneInString(stringData[i:])
					if size == 1 {
						continue
					}
				}
				v = append(v, r)
			}
			stringData = string(v)
		}
		if err := c.WriteMessage(1, []byte(stringData)); err != nil {
			log.Println(err)
			c.WriteMessage(1, errorResponse)
		}

	}()

	for {
		// c.CloseHandler()
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		} else {
			inout <- msg
		}
	}

	statusCh, errChan := dockerClient.ContainerWait(timeoutContext, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errChan:
		log.Println(err)
		c.WriteMessage(1, errorResponse)
		return
	case <-statusCh:
	}

	if err := dockerClient.ContainerRemove(timeoutContext, resp.ID, types.ContainerRemoveOptions{}); err != nil {
		log.Println(err)
	}
}
