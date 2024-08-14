package phone

import (
	"net/http"
	"peec/database"
	"peec/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PhoneNumber struct {
	Id                uint       `json:"id"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
	MobilePhoneNumber string     `json:"mobile_phone_number"`
	IsUrgency         bool       `json:"is_urgency"`
}

type UserPhoneNumber struct {
	UserId        uint `json:"user_id"`
	PhoneNumberId uint `json:"phone_number_id"`
}

/*

	ROUTES Handlers

*/

/*
ADD NEW PHONE NUMBER TO A USER BY PORVIDING user.id IN THE BODY
*/
func NewPhoneNumber(ctx *gin.Context) {
	var (
		userPhoneNumber UserPhoneNumber
		newPhone        PhoneNumber
		err             error
	)

	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Failed to retrieve params",
		})
	}

	err = ctx.ShouldBindJSON(&newPhone)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Failed to get phone number form body",
		})
		return
	}

	if !utils.IsValidPhone(newPhone.MobilePhoneNumber) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Invalid format of phone number",
		})
		return
	}

	newPhone.Id, err = database.InsertOne(newPhone)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Phone number already exist in the database",
		})
		return
	}
	// Link phone to user.
	userPhoneNumber.UserId = uint(userId)
	userPhoneNumber.PhoneNumberId = newPhone.Id
	_, err = database.InsertOne(userPhoneNumber)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Failed to link phone number to user",
		})
		return
	}

	ctx.JSON(http.StatusOK, newPhone)
	return
}

/*
UPDATE PHONE NUMBER OF A USER BY PORVIDING user.id in params AND LIST OF USER PHONE NUMBER
*/
func UpdateUserPhoneNumber(ctx *gin.Context) {
	var (
		newPhone PhoneNumber
		err      error
	)

	err = ctx.ShouldBindJSON(&newPhone)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
		return
	}
	if newPhone.Id == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "phone id is required for the operation",
		})
		return
	}
	if !utils.IsValidPhone(newPhone.MobilePhoneNumber) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Invalid format of phone number",
		})
		return
	}
	err = database.Update(newPhone)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Failed to update mobile phone number",
		})
		return
	}

	ctx.JSON(http.StatusOK, newPhone)
	return
}

/*
GET USER PHONE NUMBER BASED ON user_id PROVIDED IN PARAMS
*/

func GetUserPhoneNumber(ctx *gin.Context) {
	var (
		userId int
		phone  PhoneNumber
		err    error
	)

	userId, err = strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
	}

	err = database.Get(&phone, `SELECT phone_number.mobile_phone_number
	FROM phone_number JOIN user_phone_number 
	ON phone_number.id = user_phone_number.id 
	WHERE user_phone_number.user_id = ?`, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Failed to fetch phone number  user_id unknown",
		})
		return
	}

	ctx.JSON(http.StatusOK, phone)
}

/*
UTILS
*/
