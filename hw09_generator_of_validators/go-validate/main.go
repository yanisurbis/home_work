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
	res := generateMultipleStructValidations(parseAST())

	writeToFile(res, "models/models_validation_generated.go")

}
