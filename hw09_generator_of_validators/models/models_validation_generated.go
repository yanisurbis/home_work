
package models
import "regexp"


type ValidationError struct {
	Field string
	Err string
}

func (x User) Validate()  ([]ValidationError, error) {
errs := make([]ValidationError, 0)
	
if len(x.ID) > 36 {
	errs = append(errs, ValidationError{Field: "ID", Err: "Should be less than 36"})
}

if x.Age < 18 {
	errs = append(errs, ValidationError{Field: "Age", Err: "Should be more than 18"})
}

if x.Age > 50 {
	errs = append(errs, ValidationError{Field: "Age", Err: "Should be less than 50"})
}

{
	match, _ := regexp.MatchString("^\\w+@\\w+\\.\\w+$", x.Email)
	if !match {
		errs = append(errs, ValidationError{Field: "Email", Err: "Should satisfy the pattern ^\\w+@\\w+\\.\\w+$"})
	}
}

{
	isIn := false
	for _, v := range []string{"admin","stuff"} {
		if v == x.Role {
			isIn = true
		}
	}
	if !isIn {
		errs = append(errs, ValidationError{Field: "Role", Err: "Element should be one of "admin","stuff""})
	}
}


return errs, nil
}

func (x App) Validate()  ([]ValidationError, error) {
errs := make([]ValidationError, 0)
	
if len(x.Version) > 5 {
	errs = append(errs, ValidationError{Field: "Version", Err: "Should be less than 5"})
}


return errs, nil
}
