package main

import (
	"strings"
)

//import "strings"

//func appendErrorStr(fieldName string, )

//type FieldDescription struct {
//	Name string
//	Type string
//	Validations []FieldValidation
//}
//
//type InterfaceDescription struct {
//	Name string
//	Fields []FieldDescription
//}

func generateMultipleStructValidations(structures []InterfaceDescription) string {
	validations := `
package models
import "regexp"


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
				validationContent += generateFieldValidation(field.Name, fieldValidation)
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

func generateFieldValidation(fieldName string, validation FieldValidation) string {
	//validation = FieldValidation{
	//	Type:  "min",
	//	Value: "18",
	//}
	//fieldName := "Age"
	//fieldType := "int"
	//fieldTag := "min:18|max:50"

	validationString := ""

	if validation.Type == "min" {
		value := validation.Value.(string)
		validationString += `
if x.` + fieldName + ` < ` + value + ` {
	errs = append(errs, ValidationError{Field: "` + fieldName + `", Err: "Should be more than ` + value + `"})
}
`
	} else if validation.Type == "max" {
		value := validation.Value.(string)
		validationString += `
if x.` + fieldName + ` > ` + value + ` {
	errs = append(errs, ValidationError{Field: "` + fieldName + `", Err: "Should be less than ` + value + `"})
}
`
	} else if validation.Type == "len" {
		value := validation.Value.(string)
		validationString += `
if len(x.` + fieldName + `) > ` + value + ` {
	errs = append(errs, ValidationError{Field: "` + fieldName + `", Err: "Should be less than ` + value + `"})
}
`
	} else if validation.Type == "regexp" {
		value := validation.Value.(string)
		validationString += `
{
	match, _ := regexp.MatchString("` + value + `", x.` + fieldName + `)
	if !match {
		errs = append(errs, ValidationError{Field: "` + fieldName + `", Err: "Should satisfy the pattern ` + value + `"})
	}
}
`
	} else if validation.Type == "in" {
		values := validation.Value.([]string)
		validationString += `
{
	isIn := false
	for _, v := range []int{` + strings.Join(values, ",") + `} {
		if v == x.` + fieldName + ` {
			isIn = true
		}
	}
	if !isIn {
		errs = append(errs, ValidationError{Field: "` + fieldName + `", Err: "Element should be one of ` + strings.Join(values, ",") + `"})
	}
}
`
	}

	//fieldRules := strings.Split(fieldTag, "|")
	//for _, fieldRule := range fieldRules {
	//	ruleNameAndValue := strings.Split(fieldRule, ":")
	//	ruleName := ruleNameAndValue[0]
	//	ruleValue := ruleNameAndValue[1]
	//
	//
	//}

	return validationString
}
