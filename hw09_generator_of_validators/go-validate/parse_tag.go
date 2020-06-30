package main

import (
	"regexp"
	"strings"
)

type FieldValidation struct {
	Type  string
	Value interface{}
}

func dropFirstAndLastValue(tag string) string {
	tagLength := len(tag)
	return tag[1 : tagLength-1]
}

func removeDoubleQuotes(s string) string {
	return s[:len(s)-1]
}

func getValidations(actionValue string) []FieldValidation {
	var fieldValidations []FieldValidation

	for _, rule := range strings.Split(actionValue, "|") {
		ruleTypeAndValue := strings.Split(rule, ":")
		ruleType := ruleTypeAndValue[0]
		ruleValue := ruleTypeAndValue[1]

		// ruleType min
		// ruleValue 1
		if ruleType == "in" {
			fieldValidations = append(fieldValidations, FieldValidation{
				Type:  ruleType,
				Value: strings.Split(ruleValue, ","),
			})
		} else {
			if ruleType == "regexp" {
				_ = regexp.MustCompile(ruleValue)
			}

			fieldValidations = append(fieldValidations, FieldValidation{
				Type:  ruleType,
				Value: ruleValue,
			})
		}
	}

	return fieldValidations
}

func parseTag(fullTag string) []FieldValidation {
	tag := dropFirstAndLastValue(fullTag)

	var fieldValidations []FieldValidation

	// tag `json:"id" validate:"min:1|max:2"`
	for _, action := range strings.Split(tag, " ") {
		if len(action) > 0 {
			actionDefinition := strings.Split(action, ":\"")
			actionType := actionDefinition[0]

			// actionType validate
			// actionValue min:1|max:2
			if actionType == "validate" {
				actionValue := removeDoubleQuotes(actionDefinition[1])
				fieldValidations = append(fieldValidations, getValidations(actionValue)...)
			}
		}
	}

	return fieldValidations
}
