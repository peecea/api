package errx

import (
	"fmt"
)

func Lambda(err error) string {
	fmt.Println(err)
	return fmt.Sprintf("Something went wrong %v", err)
}

var UnAuthorizedError = "UnAuthorized"
var ParamsError = "parse params error"

var (
	ParseError            = "cannot parse"
	UnknownUserError      = "please sign up"
	IncorrectPassword     = "password incorrect"
	TypeError             = "invalid type of file"
	InvalidEmailError     = "invalid type of email"
	DuplicateUserError    = "user already exist"
	LinkUserError         = "cannot link user"
	DuplicateAddressError = "address already taken"
)

var (
	DbInsertError = "failed to insert data to database"
	DbGetError    = "failed to fetch data from database"
	DbDeleteError = "cannot delete data from the database"
	DbUpdateError = "failed to update data"
)

var (
	NeedPasswordError = "password must be set"
)
