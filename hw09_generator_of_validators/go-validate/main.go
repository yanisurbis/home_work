package main

import (
	"log"
	"os"
)

func writeToFile(str string, path string) {
	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}
	f.WriteString(str)
	f.Close()
}

func main() {
	filename := os.Getenv("GOFILE")
	// "models/models.go"
	res := generateValidation(extractInterfaceDescriptions(filename))
	writeToFile(res, "models/models_validation_generated.go")
}
