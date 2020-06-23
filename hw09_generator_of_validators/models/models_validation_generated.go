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

		if len(value) < 11 {
			errs = append(errs, ValidationError{Field: "Phones", Err: "Element on position " + strconv.Itoa(i) + " should the length should be more or equal than 11"})
			break
		}

		{
			isIn := false
			for _, v := range []string{"admin", "stuff"} {
				if v == value {
					isIn = true
				}
			}
			if !isIn {
				errs = append(errs, ValidationError{Field: "Phones", Err: "Element on position " + strconv.Itoa(i) + " should should be one of admin,stuff"})
				break
			}
		}

	}

	for i, value := range x.Phones {

		if len(value) < 11 {
			errs = append(errs, ValidationError{Field: "Phones", Err: "Element on position " + strconv.Itoa(i) + " should the length should be more or equal than 11"})
			break
		}

		{
			isIn := false
			for _, v := range []string{"admin", "stuff"} {
				if v == value {
					isIn = true
				}
			}
			if !isIn {
				errs = append(errs, ValidationError{Field: "Phones", Err: "Element on position " + strconv.Itoa(i) + " should should be one of admin,stuff"})
				break
			}
		}

	}

	return errs, nil
}
