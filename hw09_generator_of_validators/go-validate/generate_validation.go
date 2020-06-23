package main

import (
	"strings"
)

func generateMultipleStructValidations(structures []InterfaceDescription) string {
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
					validationContent += generateSliceFieldValidation(field)
				} else {
					validationContent += generateFieldValidation(field.Name, field.Type, field.TypeAlias, fieldValidation)
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

func getArrayErrorMessage(fieldName string, errorMessage string) string {
	return `errs = append(errs, ValidationError{Field: "` + fieldName + `", Err: "Element on position "+ strconv.Itoa(i) + " should ` + errorMessage + `"}) 
break`
}

// TODO: Move generateFieldValidation1 and generateFieldValidation into one function
// TODO: Handler errors during validation
// TODO: Split into separate files
func generateFieldValidation1(fieldName string, fieldType string, typeAlias string, validation FieldValidation) string {
	validationString := ""

	if validation.Type == "min" {
		value := validation.Value.(string)
		validationString += `
if value < ` + value + ` {
` + getArrayErrorMessage(fieldName, "should be more than "+value) + `	
}
`
	} else if validation.Type == "max" {
		value := validation.Value.(string)
		validationString += `
if value > ` + value + ` {
` + getArrayErrorMessage(fieldName, "should be less than "+value) + `
}
`
	} else if validation.Type == "len" {
		value := validation.Value.(string)
		validationString += `
if len(value) < ` + value + ` {
` + getArrayErrorMessage(fieldName, "the length should be more or equal than "+value) + `
}
`
	} else if validation.Type == "regexp" {
		value := validation.Value.(string)
		validationString += `
{
	match, _ := regexp.MatchString("` + value + `", value)
	if !match {
` + getArrayErrorMessage(fieldName, "should satisfy the pattern "+value) + `
	}
}
`
	} else if validation.Type == "in" {
		valuesArr := validation.Value.([]string)
		values := []string{}

		if fieldType == "[]string" {

			for _, v := range valuesArr {
				values = append(values, "\""+v+"\"")
			}
		} else {
			values = valuesArr
		}

		validationString += `
{
	isIn := false
	for _, v := range ` + typeAlias + `{` + strings.Join(values, ",") + `} {
		if v == value {
			isIn = true
		}
	}
	if !isIn {
` + getArrayErrorMessage(fieldName, "should be one of "+strings.Join(valuesArr, ",")) + `
	}
}
`
	}

	return validationString
}

func generateSliceFieldValidation(description FieldDescription) string {
	conditions := ""

	for _, validation := range description.Validations {
		conditions += generateFieldValidation1(description.Name, description.Type, description.TypeAlias, validation)
	}

	validation := `
for i, value := range x.` + description.Name + `{
	` + conditions + `
}
`

	return validation
}

func getErrorMessage(fieldName string, errorMessage string) string {
	return `errs = append(errs, ValidationError{Field: "` + fieldName + `", Err: "` + errorMessage + `"})`
}

func generateFieldValidation(fieldName string, fieldType string, typeAlias string, validation FieldValidation) string {
	validationString := ""

	// TODO use switch
	if validation.Type == "min" {
		value := validation.Value.(string)
		validationString += `
if x.` + fieldName + ` < ` + value + ` {
` + getErrorMessage(fieldName, "Should be more than "+value) + `	
}
`
	} else if validation.Type == "max" {
		value := validation.Value.(string)
		validationString += `
if x.` + fieldName + ` > ` + value + ` {
` + getErrorMessage(fieldName, "Should be less than "+value) + `
}
`
	} else if validation.Type == "len" {
		value := validation.Value.(string)
		validationString += `
if len(x.` + fieldName + `) < ` + value + ` {
` + getErrorMessage(fieldName, "The length should be more or equal than "+value) + `
}
`
	} else if validation.Type == "regexp" {
		value := validation.Value.(string)
		validationString += `
{
	match, _ := regexp.MatchString("` + value + `", x.` + fieldName + `)
	if !match {
` + getErrorMessage(fieldName, "Should satisfy the pattern "+value) + `
	}
}
`
	} else if validation.Type == "in" {
		valuesArr := validation.Value.([]string)
		values := []string{}

		if fieldType == "string" {

			for _, v := range valuesArr {
				values = append(values, "\""+v+"\"")
			}
		} else {
			values = valuesArr
		}

		validationString += `
{
	isIn := false
	for _, v := range []` + typeAlias + `{` + strings.Join(values, ",") + `} {
		if v == x.` + fieldName + ` {
			isIn = true
		}
	}
	if !isIn {
` + getErrorMessage(fieldName, "Element should be one of "+strings.Join(valuesArr, ",")) + `
	}
}
`
	}

	return validationString
}
