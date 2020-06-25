package main

import "strings"

func generatePrimitiveValidation(fieldName string, fieldType string, typeAlias string, validation FieldValidation) string {
	validationString := ""

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

func getErrorMessage(fieldName string, errorMessage string) string {
	return `errs = append(errs, ValidationError{Field: "` + fieldName + `", Err: "` + errorMessage + `"})`
}

