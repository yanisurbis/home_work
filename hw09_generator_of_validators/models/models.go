package models

type UserRole string

// NOTE: Several struct specs in one type declaration are allowed
//go:generate go-validate
//type User struct {
//	ID     string `json:"id" validate:"len:36"`
//	Name   string
//	Age    int      `validate:"min:18|max:50"`
//	Ages   []int `validate:"min:18|max:50"`
//	Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
//	Role   UserRole `validate:"in:admin,stuff"`
//	Phones []string `validate:"len:11"`
//}

type (
	User struct {
		//ID     string `json:"id" validate:"len:36"`
		//Name   string
		//Age    int      `validate:"min:18|max:50"`
		//Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		//Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11|in:admin,stuff"`
	}

	//App struct {
	//	Version string `validate:"len:5"`
	//}
)

//type Token struct {
//	Header    []byte
//	Payload   []byte
//	Signature []byte
//}
//
//type Response struct {
//	Code int    `validate:"in:200,404,500"`
//	Body string `json:"omitempty"`
//}
