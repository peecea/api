package planning

import (
	"duval/database"
	"duval/internal/authentication"
	"duval/internal/utils"
	"duval/internal/utils/errx"
	"duval/internal/utils/state"
	"duval/pkg/user"
	"duval/pkg/user/authorization"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type CalendarPlanning struct {
	Id              uint       `json:"id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	AuthorizationId uint       `json:"authorization_id"`
	StartDateTime   time.Time  `json:"start_date_time"`
	EndDateTime     time.Time  `json:"end_date_time"`
	Description     string     `json:"description"`
}

type CalendarPlanningActor struct {
	Id                 uint       `json:"id"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at"`
	AuthorizationId    uint       `json:"authorization_id"`
	CalendarPlanningId uint       `json:"calendar_planning_id"`
}

func CreateUserPlannings(ctx *gin.Context) {
	var (
		tok                   *authentication.Token
		err                   error
		calendarPlanning      CalendarPlanning
		calendarPlanningActor CalendarPlanningActor
	)

	err = ctx.ShouldBindJSON(&calendarPlanning)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		return
	}

	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		return
	}

	calendarPlanning.AuthorizationId, err = GetUserAuthorizationId(tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		return
	}

	calendarId, err := database.InsertOne(calendarPlanning)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbInsertError,
		})
		return
	}

	calendarPlanningActor.AuthorizationId = calendarPlanning.AuthorizationId
	calendarPlanningActor.CalendarPlanningId = calendarId

	err = AddCalendarPlanningActor(calendarPlanningActor)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbInsertError,
		})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, calendarPlanning)
}

func GetUserPlannings(ctx *gin.Context) {
	var (
		tok              *authentication.Token
		err              error
		authorizationId  uint
		calendarPlanning CalendarPlanning
	)

	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		return
	}

	authorizationId, err = GetUserAuthorizationId(tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		return
	}

	calendarPlanning, err = GetPlanningById(authorizationId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, calendarPlanning)
}

func RemoveUserPlannings(ctx *gin.Context) {
	var (
		tok              *authentication.Token
		err              error
		authorizationId  uint
		calendarPlanning CalendarPlanning
	)

	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		return
	}

	authorizationId, err = GetUserAuthorizationId(tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		return
	}

	calendarPlanning, err = GetPlanningById(authorizationId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		return
	}

	err = database.Delete(calendarPlanning)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbDeleteError,
		})
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
}

func AddUserIntoPlanning(ctx *gin.Context) {
	var (
		tok                   *authentication.Token
		calendarPlanningActor CalendarPlanningActor
		selectedUser          user.User
		err                   error
	)

	calendarId, err := strconv.Atoi(ctx.Param("calendar_id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.ParseError,
		})
		return
	}

	err = ctx.ShouldBindJSON(&selectedUser)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.ParseError,
		})
		return
	}

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

	authorizationId, err := GetUserAuthorizationId(selectedUser.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		return
	}

	calendarPlanningActor.AuthorizationId = authorizationId
	calendarPlanningActor.CalendarPlanningId = uint(calendarId)

	err = AddCalendarPlanningActor(calendarPlanningActor)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbInsertError,
		})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, calendarPlanningActor)
}

func GetPlanningActors(ctx *gin.Context) {
	var (
		err                    error
		calendarPlanningActors []user.User
	)

	calendarId, err := strconv.Atoi(ctx.Param("calendar_id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.ParseError,
		})
		return
	}

	calendarPlanningActors, err = GetPlanningActorByCalendarId(uint(calendarId))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbGetError,
		})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, calendarPlanningActors)
}

func RemoveUserFromPlanning(ctx *gin.Context) {
	var (
		selectedCalendarPlanningActor CalendarPlanningActor
		selectedUser                  user.User
		tok                           *authentication.Token
		err                           error
	)

	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.ParseError,
		})
		return
	}

	if tok.UserId == state.ZERO {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	calendarPlanningId, err := strconv.Atoi(ctx.Param("calendar_id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.ParamsError,
		})
		return
	}

	err = ctx.ShouldBindJSON(&selectedUser)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.ParseError,
		})
		return
	}

	selectedCalendarPlanningActor, err = GetSelectedPlanningActor(selectedUser.Id, uint(calendarPlanningId))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbGetError,
		})
		return
	}

	err = RemoveSelectedPlanningActor(selectedCalendarPlanningActor)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbDeleteError,
		})
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}

/*
	UTILS
*/

func GetUserAuthorizationId(userId uint) (id uint, err error) {
	var userAuthorization authorization.Authorization
	err = database.Get(&userAuthorization, `SELECT * FROM authorization WHERE authorization.user_id = ?`, userId)
	if err != nil {
		return 0, err
	}
	return userAuthorization.Id, nil
}

func GetPlanningById(authorizationId uint) (calendarPlanning CalendarPlanning, err error) {
	err = database.Get(&calendarPlanning, `SELECT *  FROM calendar_planning WHERE calendar_planning.authorization_id = ?`, authorizationId)
	if err != nil {
		return calendarPlanning, err
	}
	return calendarPlanning, err
}

func AddCalendarPlanningActor(calendarPlanningActor CalendarPlanningActor) (err error) {
	_, err = database.InsertOne(calendarPlanningActor)
	if err != nil {
		return err
	}
	return nil
}

func GetPlanningActorByCalendarId(calendarId uint) (calendarPlanningActors []user.User, err error) {
	err = database.GetMany(&calendarPlanningActors,
		`SELECT user.* FROM user
              JOIN authorization ON user.id = authorization.user_id
              JOIN calendar_planning_actor ON authorization.id = calendar_planning_actor.authorization_id
              JOIN calendar_planning ON calendar_planning_actor.calendar_planning_id = calendar_planning.id
     WHERE calendar_planning.id = ?`, calendarId)
	if err != nil {
		return calendarPlanningActors, err
	}
	return calendarPlanningActors, err
}

func GetSelectedPlanningActor(userId uint, calendarPlanningId uint) (calendarPlanningActor CalendarPlanningActor, err error) {
	err = database.Get(&calendarPlanningActor,
		`SELECT calendar_planning_actor.*  FROM calendar_planning_actor
                                  JOIN authorization ON calendar_planning_actor.authorization_id = authorization.id
                                  JOIN calendar_planning ON calendar_planning_actor.calendar_planning_id = calendar_planning.id
                                  JOIN user ON authorization.user_id = user.id
                                  WHERE user.id= ? AND calendar_planning.id = ?`, userId, calendarPlanningId)
	if err != nil {
		return calendarPlanningActor, err
	}
	return calendarPlanningActor, nil
}

func RemoveSelectedPlanningActor(calendarPlanningActor CalendarPlanningActor) (err error) {
	err = database.Delete(calendarPlanningActor)
	if err != nil {
		return err
	}
	return nil
}
