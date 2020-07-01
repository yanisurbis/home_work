package main

import "strings"

func generatePrimitiveFieldValidation(field FieldDescription, validation FieldValidation) string {
	fieldName := field.Name
	typeAlias := field.TypeAlias

	validationString := ""

	switch validation.Type {
	case Min:
		value := validation.Value.(string)
		validationString += `
if x.` + fieldName + ` < ` + value + ` {
` + generateErrorMessage(fieldName, "Should be more than "+value) + `	
}
`
	case Max:
		value := validation.Value.(string)
		validationString += `
if x.` + fieldName + ` > ` + value + ` {
` + generateErrorMessage(fieldName, "Should be less than "+value) + `
}
`
	case Len:
		value := validation.Value.(string)
		validationString += `
if len(x.` + fieldName + `) < ` + value + ` {
` + generateErrorMessage(fieldName, "The length should be more or equal than "+value) + `
}
`
	case Regexp:
		value := validation.Value.(string)
		validationString += `
{
	match, err := regexp.MatchString("` + value + `", x.` + fieldName + `)
	if err != nil {
		return errs, err
	}
	if !match {
` + generateErrorMessage(fieldName, "Should satisfy the pattern "+value) + `
	}
}
`
	case In:
		formattedValues, initialValues := formatSliceValues(field, validation)
		validationString += `
{
	isIn := false
	for _, v := range []` + typeAlias + `{` + strings.Join(formattedValues, ",") + `} {
		if v == x.` + fieldName + ` {
			isIn = true
		}
	}
	if !isIn {
` + generateErrorMessage(fieldName, "Element should be one of "+strings.Join(initialValues, ",")) + `
	}
}
`
	}

	return validationString
}

func generateErrorMessage(fieldName string, errorMessage string) string {
	return `errs = append(errs, ValidationError{Field: "` + fieldName + `", Err: "` + errorMessage + `"})`
}
