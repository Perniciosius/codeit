package utilities

import (
	"log"
	"os"
)

func Cleanup(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Fatalln(err)
	}
}
