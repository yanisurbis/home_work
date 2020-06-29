package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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

func getUnaliasedType(fieldType string, customTypes map[string]string) (string, error) {
	if correctType, ok := customTypes[fieldType]; ok {
		return correctType, nil
	}
	if fieldType == "string" || fieldType == "int" || fieldType == "[]string" || fieldType == "[]int" {
		return fieldType, nil
	}

	return "", fmt.Errorf("Incorrect type")
}

func isCorrectTag(tag *ast.BasicLit) bool {
	return tag != nil && strings.Contains(tag.Value, "validate:")
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

func getType1(field *ast.Field) (string, error) {
	{
		identSpec, ok := field.Type.(*ast.Ident)

		if ok {
			return identSpec.Name, nil
		}
	}

	{
		arraySpec, ok := field.Type.(*ast.ArrayType)
		if ok {
			identSpec, ok := arraySpec.Elt.(*ast.Ident)
			if ok {
				return "[]" + identSpec.Name, nil
			}
		}
	}

	return "", fmt.Errorf("not able to understand the type")
}

// TODO: Refactor
func extractInterfaceDescriptions() []InterfaceDescription {
	fs := token.NewFileSet()
	//os.Getenv("GOFILE")
	astData, _ := parser.ParseFile(fs, "models/models.go", nil, 0)
	customTypes := make(map[string]string)

	//file, _ := ioutil.ReadFile("models/models.go")
	//fileContent := string(file)

	interfaceDescriptions := []InterfaceDescription{}

	ast.Inspect(astData, func(x ast.Node) bool {
		// checking that node is a type declaration
		typeSpec, ok := x.(*ast.TypeSpec)

		if !ok {
			return true
		}

		// Create a dictionary of type aliases
		if customType := extractCustomType(typeSpec); customType != "" {
			customTypes[typeSpec.Name.Name] = customType
		}

		// checking that node is a struct
		structSpec, ok := typeSpec.Type.(*ast.StructType)

		if !ok {
			return true
		}

		fieldDescriptions := []FieldDescription{}

		for _, field := range structSpec.Fields.List {
			//func getF

			fieldType, err := getType1(field)
			fmt.Println(fieldType)

			if err != nil {
				break
			}

			unaliasedFieldType, err := getUnaliasedType(fieldType, customTypes)

			if err != nil || !isCorrectTag(field.Tag) {
				break
			}

			fieldDescriptions = append(fieldDescriptions, FieldDescription{
				Name:        field.Names[0].Name,
				Type:        unaliasedFieldType,
				TypeAlias:   fieldType,
				Validations: parseTag(field.Tag),
			})
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
