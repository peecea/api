package mark

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"peec/database"
	"peec/internal/authentication"
	"peec/internal/utils"
	"peec/internal/utils/errx"
	"peec/internal/utils/state"
	"peec/pkg/user/authorization"
	"strconv"
	"time"
)

type UserMark struct {
	Id            uint       `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
	UserId        uint       `json:"user_id"`
	AuthorId      uint       `json:"author_id"`
	AuthorComment string     `json:"author_comment"`
	AuthorMark    uint       `json:"author_mark"`
}

func RateUser(ctx *gin.Context) {
	var (
		tok         *authentication.Token
		studentMark UserMark
		err         error
	)
	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	if authorization.IsUserStudent(tok.UserId) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	if authorization.IsUserParent(tok.UserId) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	err = ctx.ShouldBindJSON(&studentMark)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.ParseError,
		})
		return
	}

	if studentMark.AuthorMark > 5 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(errors.New("value exceed 5 star")),
		})
		return
	}

	studentMark.AuthorId = tok.UserId
	err = SetUserMark(studentMark)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbInsertError,
		})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, studentMark)
}

func GetUserAverageMark(ctx *gin.Context) {
	var (
		userMarks []UserMark
		err       error
		totalMark uint
	)

	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.ParseError,
		})
		return
	}

	err = database.GetMany(&userMarks, `SELECT user_mark.* FROM user_mark WHERE user_id = ?`, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbGetError,
		})
		return
	}

	totalAuthor := len(userMarks)
	if totalAuthor == state.ZERO {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(errors.New("not rated ")),
		})
		return
	}

	for _, userMark := range userMarks {
		totalMark = totalMark + userMark.AuthorMark
	}
	averageMark := totalMark / (uint(totalAuthor))
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"average_mark": averageMark,
	})
}

func GetUserMarkComment(ctx *gin.Context) {
	var (
		tok  *authentication.Token
		err  error
		mark []UserMark
	)
	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	if authorization.IsUserStudent(tok.UserId) || authorization.IsUserParent(tok.UserId) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	err = database.GetMany(&mark,
		`SELECT user_mark.* 
			FROM user_mark
			WHERE user_mark.author_id= ?;`, tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbGetError,
		})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, mark)
}

//
/*
	UTILS
*/

func SetUserMark(userMark UserMark) (err error) {
	_, err = database.InsertOne(userMark)
	if err != nil {
		return err
	}
	return nil
}
