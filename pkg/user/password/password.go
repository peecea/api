package password

import (
	"duval/database"
	"duval/internal/utils"
	"duval/internal/utils/state"
	"errors"
	"strings"
	"time"
)

type Password struct {
	Id          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	UserId      uint       `json:"user_id"`
	Psw         string     `json:"psw"`
	ContentHash string     `json:"hash"`
}

func CreatePassword(userId uint, password Password) (err error) {
	if strings.TrimSpace(password.Psw) == state.EMPTY {
		return errors.New("password should not be empty")
	}

	if userId == state.ZERO {
		return errors.New("user id should not be empty")
	}

	if utils.PasswordHasValidLength(password.Psw) {
		return errors.New("maximum value of password is 18")
	}

	password.ContentHash, err = utils.CreateContentHash(password.Psw)
	if err != nil {
		return errors.New("password should be bound to user")
	}

	password.Psw, err = utils.HashPassword(password.Psw)
	if err != nil {
		return errors.New("something went wrong while hashing password")
	}

	password.UserId = userId
	_, err = database.Insert(password)
	if err != nil {
		return errors.New("something went wring while creating new password record on db")
	}

	return err
}

func GetUserPasswordHash(userId uint) (hash string, err error) {
	return hash, err
}

func IsPasswordValid(userId uint, passwordNaked string) bool {
	var password Password
	var err error

	err = database.Get(&password, `SELECT * FROM password WHERE user_id = ? ORDER BY created_at DESC LIMIT 1`, userId)
	if err != nil {
		return false
	}

	return utils.CheckPasswordHash(passwordNaked, password.Psw)
}
