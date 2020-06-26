
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

if x.Age < 18 {
errs = append(errs, ValidationError{Field: "Age", Err: "Should be more than 18"})	
}

if x.Age > 50 {
errs = append(errs, ValidationError{Field: "Age", Err: "Should be less than 50"})
}

{
	match, err := regexp.MatchString("^\\w+@\\w+\\.\\w+$", x.Email)
	if err != nil {
		return errs, err
	}
	if !match {
errs = append(errs, ValidationError{Field: "Email", Err: "Should satisfy the pattern ^\\w+@\\w+\\.\\w+$"})
	}
}

{
	isIn := false
	for _, v := range []UserRole{"admin","stuff"} {
		if v == x.Role {
			isIn = true
		}
	}
	if !isIn {
errs = append(errs, ValidationError{Field: "Role", Err: "Element should be one of admin,stuff"})
	}
}

for i, value := range x.Phones{
	
if len(value) < 11 {
errs = append(errs, ValidationError{Field: "Phones", Err: "Element on position "+ strconv.Itoa(i) + " should have length more or equal than 11"}) 
break
}

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
