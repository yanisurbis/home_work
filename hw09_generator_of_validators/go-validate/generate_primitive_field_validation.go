package main

import "strings"

func generatePrimitiveFieldValidation(fieldName string, fieldType string, typeAlias string, validation FieldValidation) string {
	validationString := ""

	if validation.Type == "min" {
		value := validation.Value.(string)
		validationString += `
if x.` + fieldName + ` < ` + value + ` {
` + generateErrorMessage(fieldName, "Should be more than "+value) + `	
}
`
	} else if validation.Type == "max" {
		value := validation.Value.(string)
		validationString += `
if x.` + fieldName + ` > ` + value + ` {
` + generateErrorMessage(fieldName, "Should be less than "+value) + `
}
`
	} else if validation.Type == "len" {
		value := validation.Value.(string)
		validationString += `
if len(x.` + fieldName + `) < ` + value + ` {
` + generateErrorMessage(fieldName, "The length should be more or equal than "+value) + `
}
`
	} else if validation.Type == "regexp" {
		value := validation.Value.(string)
		validationString += `
{
	match, _ := regexp.MatchString("` + value + `", x.` + fieldName + `)
	if !match {
` + generateErrorMessage(fieldName, "Should satisfy the pattern "+value) + `
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
` + generateErrorMessage(fieldName, "Element should be one of "+strings.Join(valuesArr, ",")) + `
	}
}
`
	}

	return validationString
}

func generateErrorMessage(fieldName string, errorMessage string) string {
	return `errs = append(errs, ValidationError{Field: "` + fieldName + `", Err: "` + errorMessage + `"})`
}

