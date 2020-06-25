package main

import (
	"strings"
)

type FieldValidation struct {
	Type string
	Value interface{}
}

func dropFirstAndLastValue(tag string) string {
	tagLength := len(tag)
	return tag[1:tagLength-1]
}

// TODO: check regexp validity
func parseTag(tag string) []FieldValidation {
	tag = dropFirstAndLastValue(tag)

	fieldValidations := []FieldValidation{}

	// tag `json:"id" validate:"min:1|max:2"`
	for _, action := range strings.Split(tag, " ") {
		if len(action) > 0 {
			actionDefinition := strings.Split(action, ":\"")
			actionName := actionDefinition[0]

			// actionName validate
			// actionValue min:1|max:2
			if actionName == "validate" {
				actionValue := actionDefinition[1]
				actionValueLength := len(actionValue)
				// remove double quote
				actionValue = actionValue[:actionValueLength-1]

				for _, rule := range strings.Split(actionValue, "|") {
					ruleTypeAndValue := strings.Split(rule, ":")
					ruleType := ruleTypeAndValue[0]
					ruleValue := ruleTypeAndValue[1]

					// ruleType min
					// ruleValue 1
					if ruleType == "in" {
						fieldValidations = append(fieldValidations, FieldValidation{
							Type: ruleType,
							Value: strings.Split(ruleValue, ","),
						})
					} else {
						fieldValidations = append(fieldValidations, FieldValidation{
							Type: ruleType,
							Value: ruleValue,
						})
					}
				}
			}
		}
	}

	return fieldValidations
}