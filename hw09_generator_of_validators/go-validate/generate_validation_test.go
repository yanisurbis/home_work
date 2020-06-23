package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTag(t *testing.T) {
	tags := []string{
		`json:"id" validate:"len:36"`,
		`validate:"min:18|max:50"`,
		`validate:"in:admin,stuff"`,
		`validate:"len:11"`,
		``,
		`json:"omitempty"`,
	}

	results := [][]FieldValidation{
		[]FieldValidation{FieldValidation{
			Type: "len",
			Value: "36",
		}},
		[]FieldValidation{
			FieldValidation{
				Type: "min",
				Value: "18",
			},
			FieldValidation{
				Type: "max",
				Value: "50",
			},
		},
		[]FieldValidation{FieldValidation{
			Type: "in",
			Value: []string{"admin", "stuff"},
		}},
		[]FieldValidation{FieldValidation{
			Type: "len",
			Value: "11",
		}},
		[]FieldValidation{},
		[]FieldValidation{},
	}

	t.Run("sample", func(t *testing.T) {
		for i, _ := range tags {
			require.Equal(t, results[i], parseTag(tags[i]))
		}
	})
}
