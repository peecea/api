package code

import (
	"duval/database"
	"math/rand"
	"time"
)

type Code struct {
	Id               uint       `json:"id"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
	UserId           uint       `json:"user_id"`
	VerificationCode int        `json:"value"`
}

func NewUserVerificationCode(userId uint) (err error) {
	var code Code
	code.VerificationCode = rand.Intn(9999)
	code.UserId = userId

	_, err = database.InsertOne(code)
	if err != nil {
		return err
	}

	return err
}

func IsUserVerificationCodeValid(userId uint, verificationCode int) (err error) {
	var code Code
	var query = `SELECT * FROM code WHERE user_id = ? AND verification_code = ? ORDER BY created_at desc LIMIT 1`
	err = database.Get(&code, query, userId, verificationCode)
	if err != nil {
		return err
	}

	return err
}
