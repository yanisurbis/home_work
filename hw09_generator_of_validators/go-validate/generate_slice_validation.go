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

func formatSliceValues(field FieldDescription, validation FieldValidation) ([]string, []string) {
	initialValues := validation.Value.([]string)
	var formattedValues []string

	if field.Type == StringArray || field.Type == String {
		for _, v := range initialValues {
			formattedValues = append(formattedValues, "\""+v+"\"")
		}
	} else {
		formattedValues = initialValues
	}

	return formattedValues, initialValues
}

func generateSliceElementValidation(field FieldDescription, validation FieldValidation) string {
	fieldName := field.Name
	typeAlias := field.TypeAlias

	validationString := ""

	if validation.Type == Min {
		value := validation.Value.(string)
		validationString += `
if value < ` + value + ` {
` + generateErrorForSliceElement(fieldName, "should be more than "+value) + `	
}
`
	} else if validation.Type == Max {
		value := validation.Value.(string)
		validationString += `
if value > ` + value + ` {
` + generateErrorForSliceElement(fieldName, "should be less than "+value) + `
}
`
	} else if validation.Type == Len {
		value := validation.Value.(string)
		validationString += `
if len(value) < ` + value + ` {
` + generateErrorForSliceElement(fieldName, "should have length more or equal than "+value) + `
}
`
	} else if validation.Type == Regexp {
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
	} else if validation.Type == In {
		formattedValues, initialValues := formatSliceValues(field, validation)
		validationString += `
{
	isIn := false
	for _, v := range ` + typeAlias + `{` + strings.Join(formattedValues, ",") + `} {
		if v == value {
			isIn = true
		}
	}
	if !isIn {
` + generateErrorForSliceElement(fieldName, "should be one of "+strings.Join(initialValues, ",")) + `
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
