package address

import (
	"duval/database"
	"duval/internal/authentication"
	"duval/internal/utils"
	"duval/internal/utils/errx"
	"duval/internal/utils/state"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joinverse/xid"
)

type Address struct {
	Id          uint       `json:"id"`
	Country     string     `json:"country"`
	City        string     `json:"city"`
	Latitude    float64    `json:"latitude"`
	Longitude   float64    `json:"longitude"`
	Street      string     `json:"street"`
	FullAddress string     `json:"full_address"`
	Xid         string     `json:"xid"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type UserAddress struct {
	Id          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	UserId      uint       `json:"user_id"`
	AddressId   uint       `json:"address_id"`
	AddressType string     `json:"address_type"`
}

func NewAddress(ctx *gin.Context) {

	var (
		tok         *authentication.Token
		userId      uint
		address     Address
		userAddress UserAddress
		err         error
	)
	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}
	userId = uint(tok.UserId)
	err = ctx.ShouldBindJSON(&address)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.ParseError,
		})
		return
	}
	// get user address
	isUser, err := GetUserAddressWithId(userId)
	if isUser.AddressId > state.ZERO {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DuplicateAddressError,
		})
		return
	}
	address.Xid = xid.New().String()

	address.Id, err = database.InsertOne(address)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbInsertError,
		})
		return
	}

	// Link new address to the current user
	userAddress.UserId = userId
	userAddress.AddressId = address.Id
	_, err = database.InsertOne(userAddress)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.LinkUserError,
		})
		return
	}

	ctx.JSON(http.StatusOK, address)
	return
}

/*
UPDATE ADDRESS OF A USER BY PROVIDING ID IN THE BODY
*/
func UpdateUserAddress(ctx *gin.Context) {
	var (
		address Address
		err     error
	)

	err = ctx.ShouldBindJSON(&address)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
		return
	}
	if address.Id == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "address id is required for the operation",
		})
		return
	}

	err = database.Update(address)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbUpdateError,
		})
		return
	}

	ctx.JSON(http.StatusOK, address)
	return
}

/*
GET USER ADDRESS  BASED ON user_id PROVIDED IN PARAMS
*/
func GetUserAddress(ctx *gin.Context) {
	var (
		tok *authentication.Token

		userId  uint
		address Address
		err     error
	)

	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	userId = uint(tok.UserId)

	err = database.Get(&address, `SELECT address.*
    FROM address JOIN user_address 
    ON address.id = user_address.address_id 
    WHERE user_address.user_id = ?`, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbGetError,
		})
		return
	}

	ctx.JSON(http.StatusOK, address)
}

/*
REMOVE USER ADDRESS  BASED ON user_id PROVIDED IN PARAMS
*/

func RemoveUserAddress(ctx *gin.Context) {
	var (
		tok *authentication.Token

		userId      uint
		address     Address
		userAddress UserAddress
		err         error
	)

	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}
	userId = uint(tok.UserId)

	err = database.Get(&address, `SELECT address.*
    FROM address JOIN user_address 
    ON address.id = user_address.address_id 
    WHERE user_address.user_id = ?`, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbGetError,
		})
		return
	}

	err = database.Delete(address)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbDeleteError,
		})
		return
	}
	// and remove user_address

	err = database.Get(&userAddress, `SELECT * FROM user_address where user_id = ?`, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbGetError,
		})
		return
	}
	err = database.Delete(userAddress)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbDeleteError,
		})
		return
	}
	ctx.JSON(http.StatusOK, "Address removed successfuly!")
}

/*
GET USER_ADDRESS WITH USER_ID
*/

func GetUserAddressWithId(userId uint) (userAddress UserAddress, err error) {
	err = database.Get(&userAddress, "SELECT * FROM user_address Where user_id = ?", userId)
	if err != nil {
		return userAddress, err
	}
	return userAddress, err
}
