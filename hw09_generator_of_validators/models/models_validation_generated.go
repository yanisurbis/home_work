
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
	
if len(x.ID) < 36 {
errs = append(errs, ValidationError{Field: "ID", Err: "The length should be more or equal than 36"})
}


return errs, nil
}

func (x App) Validate()  ([]ValidationError, error) {
errs := make([]ValidationError, 0)
	
if len(x.Version) < 5 {
errs = append(errs, ValidationError{Field: "Version", Err: "The length should be more or equal than 5"})
}


return errs, nil
}

func (x Response) Validate()  ([]ValidationError, error) {
errs := make([]ValidationError, 0)
	
{
	isIn := false
	for _, v := range []int{200,404,500} {
		if v == x.Code {
			isIn = true
		}
	}
	if !isIn {
errs = append(errs, ValidationError{Field: "Code", Err: "Element should be one of 200,404,500"})
	}
}


return errs, nil
}
