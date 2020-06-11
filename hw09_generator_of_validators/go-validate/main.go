package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

//func generateValidation(fieldName string, fieldType string, fieldTag string) {
//func generateFieldValidation(fieldName string, fieldType string, fieldTag string) {
//	fieldName := "Age"
//	fieldType := "int"
//	fieldTag := "min:18|max:50"
//}

func generateStructValidation() string {
	return `
// Code generated by cool go-validate tool; DO NOT EDIT.
package models

func (u User) Validate() ([]ValidationError, error) {
    errs := make([]ValidationError, 1)
	
	if (true) {
		errs = append(errs, &ValidationError{Field: "Age", Err: "Should be less than 50"}
	}

	return errs, nil
}
`
}

func parseAST() {
	fs := token.NewFileSet()
	//os.Getenv("GOFILE")
	astData, _ := parser.ParseFile(fs, "models/models.go", nil, 0)
	//println(astData)
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
			//fmt.Println(field.Names)
			fmt.Println(field.Tag)
		}

		return false
	})
}

func main() {
//	f, _ := os.Create("models/models_validation_generated.go")
//	f.WriteString(`
//// Code generated by cool go-validate tool; DO NOT EDIT.
//package models
//
//func (u User) Validate() ([]ValidationError, error) {
//	return nil, nil
//}
//`)
//	f.Close()
	parseAST()
}
