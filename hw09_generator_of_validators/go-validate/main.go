package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

//func generateValidation(fieldName string, fieldType string, fieldTag string) {
func generateFieldValidation() string {
	fieldName := "Age"
	fieldType := "int"
	fieldTag := "min:18|max:50"

	validation := ""

	fieldRules := strings.Split(fieldTag, "|")
	for _, fieldRule := range fieldRules {
		ruleNameAndValue := strings.Split(fieldRule, ":")
		ruleName := ruleNameAndValue[0]
		ruleValue := ruleNameAndValue[1]

		if (fieldType == "int") {
			if (ruleName == "min") {
				validation += `
if (x.` + fieldName + ` < ` + ruleValue + `) {
	errs = append(errs, ValidationError{Field: "` + fieldName + `", Err: "Should be more than ` + ruleValue + `"})
}
`
			} else if (ruleName == "max") {
				validation += `
if (x.` + fieldName + ` > ` + ruleValue + `) {
	errs = append(errs, ValidationError{Field: "` + fieldName + `", Err: "Should be less than ` + ruleValue + `"})
}
`
			}
		}
	}

	return validation
}

func generateStructValidation() string {
	return `
// Code generated by cool go-validate tool; DO NOT EDIT.
package models

func (x User) Validate() ([]ValidationError, error) {
    errs := make([]ValidationError, 0) `+ generateFieldValidation() +`

	return errs, nil
}
`
}

func getType(content string, start int, end int) string {
	return content[start:end]
}

func parseAST() {
	fs := token.NewFileSet()
	//os.Getenv("GOFILE")
	astData, _ := parser.ParseFile(fs, "models/models.go", nil, 0)
	//println(astData)

	file, _ := ioutil.ReadFile("models/models.go")
	fileContent := string(file)

	ast.Inspect(astData, func(x ast.Node) bool {
		typeSpec, ok := x.(*ast.TypeSpec)

		if !ok {
			return true
		}

		structSpec, ok := typeSpec.Type.(*ast.StructType)

		if !ok {
			return true
		}

		fmt.Println(typeSpec.Name)

		for _, field := range structSpec.Fields.List {
			//fmt.Println(field.Type)
			//fmt.Println(field.Type.End())
			//fmt.Println(field.Type.Pos())
			fmt.Println(getType(fileContent, int(field.Type.Pos()) - 1, int(field.Type.End()) - 1))
			fmt.Println(field.Names[0])
			if field.Tag != nil {
				fmt.Println(field.Tag.Value)
			}
		}

		fmt.Println("==================")

		return false
	})
}



func main() {
	//println(generateStructValidation())
	//f, _ := os.Create("models/models_validation_generated.go")
	//f.WriteString(generateStructValidation())
	//f.Close()
	parseAST()

	//dat, _ := ioutil.ReadFile("models/models.go")
	//fmt.Print(string(dat))
//	generateFieldValidation()
}
