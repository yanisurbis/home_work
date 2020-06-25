package main

import "strings"

func generateSliceValidation(description FieldDescription) string {
	conditions := ""

	for _, validation := range description.Validations {
		conditions += generateSliceElementValidation(description.Name, description.Type, description.TypeAlias, validation)
	}

	validation := `
for i, value := range x.` + description.Name + `{
	` + conditions + `
}
`

	return validation
}

// TODO: Handler errors during validation
func generateSliceElementValidation(fieldName string, fieldType string, typeAlias string, validation FieldValidation) string {
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

func getArrayErrorMessage(fieldName string, errorMessage string) string {
	return `errs = append(errs, ValidationError{Field: "` + fieldName + `", Err: "Element on position "+ strconv.Itoa(i) + " should ` + errorMessage + `"}) 
break`
}