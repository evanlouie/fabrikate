package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cwd)
	p := filepath.Join("internal", "generatable", "_generated")
	files, err := os.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		log.Println(file.Name())
		fileP := filepath.Join(p, file.Name())
		b, err := os.ReadFile(fileP)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(len(b))
	}
}
