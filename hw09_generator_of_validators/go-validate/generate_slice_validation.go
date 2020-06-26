package main

import "strings"

func generateSliceValidation(field FieldDescription) string {
	conditions := ""

	for _, validation := range field.Validations {
		conditions += generateSliceElementValidation(field, validation)
	}

	validation := `
for i, value := range x.` + field.Name + `{
	` + conditions + `
}
`

	return validation
}

func generateSliceElementValidation(field FieldDescription, validation FieldValidation) string {
	fieldName := field.Name
	fieldType := field.Type
	typeAlias := field.TypeAlias

	validationString := ""

	if validation.Type == "min" {
		value := validation.Value.(string)
		validationString += `
if value < ` + value + ` {
` + generateErrorForSliceElement(fieldName, "should be more than "+value) + `	
}
`
	} else if validation.Type == "max" {
		value := validation.Value.(string)
		validationString += `
if value > ` + value + ` {
` + generateErrorForSliceElement(fieldName, "should be less than "+value) + `
}
`
	} else if validation.Type == "len" {
		value := validation.Value.(string)
		validationString += `
if len(value) < ` + value + ` {
` + generateErrorForSliceElement(fieldName, "should have length more or equal than "+value) + `
}
`
	} else if validation.Type == "regexp" {
		value := validation.Value.(string)
		validationString += `
{
	match, err := regexp.MatchString("` + value + `", value)
	if err != nil {
		return errs, err
	}
	if !match {
` + generateErrorForSliceElement(fieldName, "should satisfy the pattern "+value) + `
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
` + generateErrorForSliceElement(fieldName, "should be one of "+strings.Join(valuesArr, ",")) + `
	}
}
`
	}

	return validationString
}

func generateErrorForSliceElement(fieldName string, errorMessage string) string {
	return `errs = append(errs, ValidationError{Field: "` + fieldName + `", Err: "Element on position "+ strconv.Itoa(i) + " ` + errorMessage + `"}) 
break`
}