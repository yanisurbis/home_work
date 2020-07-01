package main

import (
	"os"
)

func writeToFile(str string, path string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	_, _ = f.WriteString(str)
	_ = f.Close()
}

func main() {
	filename := os.Getenv("GOFILE")
	res := generateValidation(extractInterfaceDescriptions(filename))
	//res := generateValidation(extractInterfaceDescriptions("models/models.go"))
	writeToFile(res, os.Getenv("PWD")+"/models/models_validation_generated.go")
}
