package models

import (
	"strconv"
)

type ValidationError struct {
	Field string
	Err   string
}

func (x User) Validate() ([]ValidationError, error) {
	errs := make([]ValidationError, 0)

	for i, value := range x.Phones {

		if value < 11 {
			errs = append(errs, ValidationError{Field: "Phones", Err: "Element on position " + strconv.Itoa(i) + " should should be more than 11"})
			break
		}

		{
			isIn := false
			for _, v := range []int{12, 13} {
				if v == value {
					isIn = true
				}
			}
			if !isIn {
				errs = append(errs, ValidationError{Field: "Phones", Err: "Element on position " + strconv.Itoa(i) + " should should be one of 12,13"})
				break
			}
		}

	}

	for i, value := range x.Phones {

		if value < 11 {
			errs = append(errs, ValidationError{Field: "Phones", Err: "Element on position " + strconv.Itoa(i) + " should should be more than 11"})
			break
		}

		{
			isIn := false
			for _, v := range []int{12, 13} {
				if v == value {
					isIn = true
				}
			}
			if !isIn {
				errs = append(errs, ValidationError{Field: "Phones", Err: "Element on position " + strconv.Itoa(i) + " should should be one of 12,13"})
				break
			}
		}

	}

	return errs, nil
}
