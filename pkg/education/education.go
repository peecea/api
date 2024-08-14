package education

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

type Education struct {
	Id        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name"`
}

type Subject struct {
	Id               uint       `json:"id"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
	EducationLevelId uint       `json:"education_level_id"`
	Name             string     `json:"name"`
}

type UserEducationLevelSubject struct {
	Id        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	UserId    uint       `json:"user_id"`
	SubjectId uint       `json:"subject_id"`
}

func GetSubjects(ctx *gin.Context) {
	var (
		err      error
		subjects []Subject
		eduId    int
	)

	eduId, err = strconv.Atoi(ctx.Param("edu"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Failed to retrieve params",
		})
	}

	err = database.Select(&subjects, `SELECT * FROM subject WHERE education_level_id = ?`, eduId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		return
	}

	ctx.JSON(http.StatusOK, subjects)
	return
}

func GetEducation(ctx *gin.Context) {
	var (
		err  error
		edus []Education
	)

	err = database.Select(&edus, `SELECT * FROM education WHERE id > 0 ORDER BY  created_at`)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		return
	}

	ctx.JSON(http.StatusOK, edus)
	return
}

// User educationLevel

func SetUserEducationLevel(ctx *gin.Context) {
	var (
		tok                       *authentication.Token
		userEducationLevelSubject UserEducationLevelSubject
		subject                   Subject
		err                       error
	)
	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}
	if tok.UserId == state.ZERO {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	if !authorization.IsUserStudent(tok.UserId) {
		if !authorization.IsUserProfessor(tok.UserId) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
				Message: errx.Lambda(errors.New("not as student or professor")),
			})
			return
		}
	}

	err = ctx.ShouldBindJSON(&subject)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.ParseError,
		})
		return
	}

	userEducationLevelSubject.SubjectId = uint(subject.Id)
	userEducationLevelSubject.UserId = tok.UserId

	_, err = database.InsertOne(userEducationLevelSubject)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbInsertError,
		})
		return
	}

	userLevel, err := GetUserLevel(tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{Message: errx.DbGetError})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, userLevel)
}

func GetUserEducationLevel(ctx *gin.Context) {
	var (
		err error
		tok *authentication.Token
	)

	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	userLevel, err := GetUserLevel(tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbGetError,
		})
		return

	}
	ctx.AbortWithStatusJSON(http.StatusOK, userLevel)

}

func UpdateUserEducationLevel(ctx *gin.Context) {
	var (
		tok                              *authentication.Token
		currentUserEducationLevelSubject UserEducationLevelSubject
		userEducationLevelSubject        UserEducationLevelSubject
		subject                          Subject
		err                              error
	)

	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}
	if tok.UserId == state.ZERO {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	if !authorization.IsUserStudent(tok.UserId) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Not a student",
		})
		return
	}

	err = ctx.ShouldBindJSON(&subject)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.ParseError,
		})
		return
	}

	err = database.Get(&currentUserEducationLevelSubject, `SELECT user_education_level_subject.* FROM user_education_level_subject
			WHERE user_education_level_subject.user_id = ?`, tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbGetError,
		})
		return
	}

	err = RemoveUserEducationLevelSubject(currentUserEducationLevelSubject)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbDeleteError,
		})
		return
	}

	userEducationLevelSubject.SubjectId = subject.Id
	userEducationLevelSubject.UserId = tok.UserId
	_, err = database.InsertOne(userEducationLevelSubject)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbInsertError,
		})
		return
	}

	userLevel, err := GetUserLevel(tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{Message: errx.DbGetError})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, userLevel)
}

// User education subjects

func GetUserSubjects(ctx *gin.Context) {
	var (
		subjects []Subject
		tok      *authentication.Token
		err      error
	)
	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	if authorization.IsUserParent(tok.UserId) || authorization.IsUserTutor(tok.UserId) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	if authorization.IsUserStudent(tok.UserId) {
		err = database.GetMany(&subjects, `SELECT subject.* 
											FROM subject
											WHERE subject.education_level_id = (SELECT education.id FROM education  JOIN subject ON education.id  =  subject.education_level_id JOIN user_education_level_subject ON subject.id = user_education_level_subject.subject_id
                                   			WHERE user_education_level_subject.user_id = ?)`, tok.UserId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
				Message: errx.DbGetError,
			})
			return
		}
	}

	if authorization.IsUserProfessor(tok.UserId) {
		err = database.GetMany(&subjects, `SELECT subject.* FROM subject
			JOIN user_education_level_subject  ON subject.id = user_education_level_subject.subject_id
			WHERE user_education_level_subject.user_id = ?`, tok.UserId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
				Message: errx.DbGetError,
			})
			return
		}
	}

	ctx.AbortWithStatusJSON(http.StatusOK, subjects)
}

/*
	UTILS
*/

func GetUserLevel(userId uint) (educationLevel Education, err error) {

	err = database.Get(&educationLevel,
		`SELECT education.* FROM education
				JOIN subject ON education.id  =  subject.education_level_id
				JOIN user_education_level_subject ON subject.id = user_education_level_subject.subject_id
			WHERE user_education_level_subject.user_id = ?`, userId)
	if err != nil {
		return educationLevel, err
	}

	return educationLevel, err
}

func RemoveUserEducationLevelSubject(userEducationLevelSubject UserEducationLevelSubject) (err error) {
	err = database.Delete(userEducationLevelSubject)
	if err != nil {
		return err
	}
	return err
}
