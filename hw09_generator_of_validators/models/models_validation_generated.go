
package models
import (
	"regexp"
	"strconv"
)

type ValidationError struct {
	Field string
	Err string
}

func (x User) Validate()  ([]ValidationError, error) {
errs := make([]ValidationError, 0)
	
if x.Phones < 11 {
errs = append(errs, ValidationError{Field: "Phones", Err: Should be more than 11})	
}

{
	isIn := false
	for _, v := range [][]int{12,13} {
		if v == x.Phones {
			isIn = true
		}
	}
	if !isIn {
errs = append(errs, ValidationError{Field: "Phones", Err: Element should be one of 12,13})
	}
}


return errs, nil
}
