package main

import (
	"os"
)

func writeToFile(str string, path string) {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	f.WriteString(str)
}

func main() {
	filename := os.Getenv("GOFILE")
	res := generateValidation(extractInterfaceDescriptions(filename))
	writeToFile(res, os.Getenv("PWD") + "/models/models_validation_generated.go")
}
