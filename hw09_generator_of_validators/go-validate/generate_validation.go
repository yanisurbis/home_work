package main

func generateValidation(structures []InterfaceDescription) string {
	validations := `
package models
import (
	"regexp"
	"strconv"
)

type ValidationError struct {
	Field string
	Err string
}
`

	for _, structure := range structures {
		validations += generateStructValidation(structure)
	}

	return validations
}

func generateStructValidation(structure InterfaceDescription) string {
	validations := func(fields []FieldDescription) string {
		validationContent := ""

		for _, field := range fields {
			for _, fieldValidation := range field.Validations {
				if field.Type == "[]string" || field.Type == "[]int" {
					validationContent += generateSliceValidation(field)
				} else {
					validationContent += generatePrimitiveFieldValidation(field, fieldValidation)
				}
			}
		}

		return validationContent
	}(structure.Fields)

	return `
func (x ` + structure.Name + `) Validate()  ([]ValidationError, error) {
errs := make([]ValidationError, 0)
	` + validations + `

return errs, nil
}
`
}
