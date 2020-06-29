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

func getType(t ast.Expr) (string, error) {
	if identSpec, ok := t.(*ast.Ident); ok {
		if ok {
			return identSpec.Name, nil
		}
	} else if arraySpec, ok := t.(*ast.ArrayType); ok {
		identSpec, ok := arraySpec.Elt.(*ast.Ident)
		if ok {
			return "[]" + identSpec.Name, nil
		}
	}

	return "", fmt.Errorf("not able to understand the type")
}

func extractInterfaceDescriptions(filename string) []InterfaceDescription {
	fs := token.NewFileSet()
	astData, _ := parser.ParseFile(fs, filename, nil, 0)

	// store type aliases
	customTypes := make(map[string]string)

	interfaceDescriptions := []InterfaceDescription{}

	ast.Inspect(astData, func(x ast.Node) bool {
		// checking that node is a type declaration
		typeSpec, ok := x.(*ast.TypeSpec)

		if !ok {
			return true
		}


		// Create a dictionary of type aliases
		if customType, err := getType(typeSpec.Type); err == nil {
			customTypes[typeSpec.Name.Name] = customType
		}

		// checking that node is a struct
		structSpec, ok := typeSpec.Type.(*ast.StructType)

		if !ok {
			return true
		}

		fieldDescriptions := []FieldDescription{}

		for _, field := range structSpec.Fields.List {
			fieldType, err := getType(field.Type)

			if err != nil {
				continue
			}

			unaliasedFieldType, err := getUnaliasedType(fieldType, customTypes)

			if err != nil || !isCorrectTag(field.Tag) {
				continue
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
