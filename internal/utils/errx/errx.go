package errx

import (
	"errors"
	"fmt"
)

func Lambda(err error) string {
	fmt.Println(err)
	return "Something went wrong"
}

var UnAuthorizedError = errors.New("UnAuthorized")
var ParamsError = errors.New("parse params error")

var (
	ParseError            = errors.New("cannot parse")
	TypeError             = errors.New("invalid type of file")
	InvalidEmailError     = errors.New("invalid type of email")
	DuplicateUserError    = errors.New("user already exist")
	LinkUserError         = errors.New("cannot link user")
	DuplicateAddressError = errors.New("address already taken")
)

var (
	DbInsertError = errors.New("failed to insert data to database")
	DbGetError    = errors.New("failed to fetch data from database")
	DbDeleteError = errors.New("cannot delete data from the database")
	DbUpdateError = errors.New("failed to update data")
)

var (
	NeedPasswordError = errors.New("password must be set")
)
