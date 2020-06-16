package main

import (
	"strings"
)

type FieldValidation struct {
	Type string
	Value interface{}
}

func parseTag(tag string) []FieldValidation {

	fieldValidations := []FieldValidation{}

	for _, action := range strings.Split(tag, " ") {
		actionDefinition := strings.Split(action, ":\"")
		actionName := actionDefinition[0]
		actionValue := actionDefinition[1]
		actionValueLength := len(actionValue)
		actionValue = actionValue[:actionValueLength-1]

		if actionName == "validate" {
			for _, rule := range strings.Split(actionValue, "|") {
				ruleTypeAndValue := strings.Split(rule, ":")
				fieldValidations = append(fieldValidations, FieldValidation{
					Type: ruleTypeAndValue[0],
					Value: ruleTypeAndValue[1],
				})
			}
		}
	}

	return fieldValidations
}