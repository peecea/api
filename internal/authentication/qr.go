package authentication

import (
	"duval/database"
	"duval/internal/configuration"
	"duval/internal/utils"
	"duval/internal/utils/errx"
	"github.com/gin-gonic/gin"
	"github.com/joinverse/xid"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"net/http"
	"strconv"
	"time"
)

type QrCodeRegistry struct {
	Id        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	UserId    uint       `json:"user_id"`
	Xid       string     `json:"xid"`
	IsUsed    bool       `json:"is_used"`
}

func GenerateQrCode(ctx *gin.Context) {
	var (
		tok            *Token
		err            error
		networkLink    string
		qrImageLink    string
		qrCodeRegistry QrCodeRegistry
	)
	time.Sleep(100)
	tok, err = GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	networkLink = "http://" + configuration.App.Host + ":" + configuration.App.Port + "/api/login/with-qr/:" + strconv.Itoa(int(tok.UserId))

	qrc, err := qrcode.New(networkLink)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
		return
	}

	qrCodeRegistry.UserId = tok.UserId
	qrCodeRegistry.Xid = xid.New().String()
	qrCodeRegistry.IsUsed = false

	qrImageLink = "http://" + configuration.App.Host + ":" + configuration.App.Port + "/api/public/qr/" + qrCodeRegistry.Xid + ".jpg"

	w, err := standard.New(utils.FILE_UPLOAD_DIR + utils.QR_CODE_UPLOAD_DIR + qrCodeRegistry.Xid + ".jpg")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
		return
	}
	// save file
	err = qrc.Save(w)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
		return
	}
	_, err = database.InsertOne(qrCodeRegistry)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbInsertError,
		})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, qrImageLink)
}

func LoginWithQr(ctx *gin.Context) {
	time.Sleep(100)
	var (
		tok            string
		err            error
		qrCodeRegistry QrCodeRegistry
	)

	xId := ctx.Param("xid")

	//Retrieve QR Code Registry: Retrieves the qr_code_registry record that corresponds to a given xid (likely a unique identifier for the QR code).
	qrCodeRegistry, err = GetQrCodeRegistry(xId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
		return
	}

	//Update is_used Flag: Sets the is_used field of the retrieved qr_code_registry to true, indicating the QR code has been used for login.
	err = UpdateQrCodeRegistryFlag(qrCodeRegistry)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbUpdateError,
		})
		return
	}

	//	Generate Access Token: Creates an access token using the user_id from the qr_code_registry to authenticate the user's session
	tok, err = GetTokenString(qrCodeRegistry.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"token": tok,
	})
}

/*
	UTILS
*/

func GetQrCodeRegistry(xId string) (qrCodeRegistry QrCodeRegistry, err error) {
	err = database.Get(&qrCodeRegistry, `SELECT * FROM qr_code_registry WHERE qr_code_registry.xid = ?`, xId)
	if err != nil {
		return qrCodeRegistry, err
	}
	return qrCodeRegistry, nil
}

func UpdateQrCodeRegistryFlag(qrCodeRegistry QrCodeRegistry) (err error) {
	qrCodeRegistry.IsUsed = true
	err = database.Update(qrCodeRegistry)
	if err != nil {
		return err
	}
	return
}
