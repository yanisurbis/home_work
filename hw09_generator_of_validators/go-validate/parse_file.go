package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

type FieldDescription struct {
	Name        string
	Type        string
	TypeAlias   string
	Validations []FieldValidation
}

type InterfaceDescription struct {
	Name   string
	Fields []FieldDescription
}

func getType(content string, start int, end int) string {
	return content[start:end]
}

func extractCustomType(typeSpec *ast.TypeSpec) string {
	// simple type: string, int
	identSpec, ok := typeSpec.Type.(*ast.Ident)

	if ok && (identSpec.Name == "string" || identSpec.Name == "int") {
		return identSpec.Name
	}

	// complex type: []string, []int
	arraySpec, ok := typeSpec.Type.(*ast.ArrayType)

	if ok {
		arrayElmSpec, ok := arraySpec.Elt.(*ast.Ident)
		if ok && (arrayElmSpec.Name == "string" || arrayElmSpec.Name == "int") {
			return "[]" + arrayElmSpec.Name
		}
	}

	return ""
}

// TODO: Refactor
func parseAST() []InterfaceDescription {
	fs := token.NewFileSet()
	//os.Getenv("GOFILE")
	astData, _ := parser.ParseFile(fs, "models/models.go", nil, 0)
	customTypes := make(map[string]string)

	file, _ := ioutil.ReadFile("models/models.go")
	fileContent := string(file)

	interfaceDescriptions := []InterfaceDescription{}

	ast.Inspect(astData, func(x ast.Node) bool {
		typeSpec, ok := x.(*ast.TypeSpec)

		if !ok {
			return true
		}

		if customType := extractCustomType(typeSpec); customType != "" {
			customTypes[typeSpec.Name.Name] = customType
		}

		structSpec, ok := typeSpec.Type.(*ast.StructType)

		if !ok {
			return true
		}

		fieldDescriptions := []FieldDescription{}

		for _, field := range structSpec.Fields.List {
			fieldType := getType(fileContent, int(field.Type.Pos())-1, int(field.Type.End())-1)

			correctFieldType := func(fieldType string, customTypes map[string]string) string {
				if correctType, ok := customTypes[fieldType]; ok {
					return correctType
				}
				if fieldType == "string" || fieldType == "int" || fieldType == "[]string" || fieldType == "[]int" {
					return fieldType
				}
				return ""
			}(fieldType, customTypes)

			isCorrectTag := field.Tag != nil && strings.Contains(field.Tag.Value, "validate:")

			if isCorrectTag && correctFieldType != "" {
				fieldDescriptions = append(fieldDescriptions, FieldDescription{
					Name:        field.Names[0].Name,
					Type:        correctFieldType,
					TypeAlias:   fieldType,
					Validations: parseTag(field.Tag.Value),
				})
			}
		}

		if len(fieldDescriptions) != 0 {
			interfaceDescriptions = append(interfaceDescriptions, InterfaceDescription{
				Name:   typeSpec.Name.Name,
				Fields: fieldDescriptions,
			})
		}

		return false
	})

	return interfaceDescriptions
}
