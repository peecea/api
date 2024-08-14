package authentication

import (
	"errors"
	"peec/database"
	"peec/internal/configuration"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	UserId     uint `json:"user_id"`
	UserLevel  uint `json:"user_level"`
	UserStatus uint `json:"user_status"`
	jwt.RegisteredClaims
}

func GetTokenString(userId uint) (str string, err error) {
	var tok Token
	err = database.Get(&tok, `SELECT u.id as 'user_id', u.status as 'user_status' FROM user u WHERE u.id = ?`, userId)
	if err != nil {
		return str, err
	}

	err = database.Get(&tok, `SELECT auth.level as 'user_level' FROM authorization auth WHERE auth.user_id = ?`, userId)
	if err != nil {
		return str, err
	}

	str, err = NewAccessToken(tok)
	if err != nil {
		return str, err
	}

	return str, err
}

func GetTokenDataFromContext(ctx *gin.Context) (tok *Token, err error) {
	tokenString := ctx.GetHeader("Authorization")
	if len(strings.TrimSpace(tokenString)) == 0 {
		return nil, errors.New("bad header value given")
	}

	bearer := strings.Split(tokenString, " ")
	if len(bearer) != 2 {
		return nil, errors.New("incorrectly formatted authorization header")
	}

	return ParseAccessToken(bearer[1]), err
}

func NewAccessToken(claims Token) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(configuration.App.TokenSecret))
}

func NewRefreshToken(claims jwt.RegisteredClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(configuration.App.TokenSecret))
}

func ParseAccessToken(accessToken string) *Token {
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &Token{}, func(token *jwt.Token) (any, error) {
		return []byte(configuration.App.TokenSecret), nil
	})

	return parsedAccessToken.Claims.(*Token)
}

func ParseRefreshToken(refreshToken string) *jwt.RegisteredClaims {
	parsedRefreshToken, _ := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(configuration.App.TokenSecret), nil
	})

	return parsedRefreshToken.Claims.(*jwt.RegisteredClaims)
}
