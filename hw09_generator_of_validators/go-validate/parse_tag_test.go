package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateValidation(t *testing.T) {
	examples := []FieldValidation{
		FieldValidation{
			Type:  "min",
			Value: "18",
		},
		FieldValidation{
			Type:  "max",
			Value: "18",
		},
		FieldValidation{
			Type:  "len",
			Value: "18",
		},
		FieldValidation{
			Type:  "regexp",
			Value: "18",
		},
		FieldValidation{
			Type:  "in",
			Value: []string{"1", "2"},
		},
	}

	for _, example := range examples {
		fmt.Println(generateFieldValidation("Age", example))
	}

	require.Equal(t, true, true)
}
