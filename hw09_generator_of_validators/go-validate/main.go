package main

import "regexp"

//func generateStructValidation() string {
//	return `
//// Code generated by cool go-validate tool; DO NOT EDIT.
//package models
//
//func (x User) Validate() ([]ValidationError, error) {
//    errs := make([]ValidationError, 0) `+ generateFieldValidation() +`
//
//	return errs, nil
//}
//`
//}

func main() {
	//println(generateStructValidation())
	//f, err := os.Create("models_validation_generated.go")
	//if err != nil {
	//	log.Println(err)
	//}
	//f.WriteString(generateStructValidation())
	//f.Close()
	//for _, v := range parseAST() {
	//	fmt.Printf("%+v\n\n\n", v)
	//}

	res := generateMultipleStructValidations(parseAST())
	//fmt.Println(res)
	regexp.MustCompile(res)
	//	regexp.MustCompile(`
	//package main
	//
	//type User struct {
	//	ID string
	//}
	//type ValidationError struct {
	//	Field string
	//	Err string
	//}
	//
	//
	//
	//func (x User) Validate()  ([]ValidationError, error) {
	//	errors := make([]ValidationError, 0)
	//
	//
	//	return errors, nil
	//}
	//`)
	//fmt.Println(os.Getenv("GOFILE"))

	//path, err := os.Getwd()
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Println(path)

	//dat, _ := ioutil.ReadFile("models/models.go")
	//fmt.Print(string(dat))
	//	generateFieldValidation()
	//	fmt.Println(parseTag(""))
	//	x := 2
	//	values := []string{"1", "2", "3"}
	//	isIn := false
	//	for _, v := range values {
	//		if v == x {
	//			isIn = true
	//		}
	//	}
	//	if isIn {
	//	//	add error
	//	}
	//	fmt.Println(strings.Join(values, ", "))
}
