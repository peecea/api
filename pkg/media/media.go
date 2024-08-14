package media

import (
	"github.com/gabriel-vasile/mimetype"
	"net/http"
	"path/filepath"
	"peec/database"
	"peec/internal/authentication"
	"peec/internal/utils"
	"peec/internal/utils/errx"
	"peec/internal/utils/state"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joinverse/xid"
)

const (
	CV               = 0
	CoverLetter      = 1
	Video            = 2
	UserProfileImage = 3
)

type Media struct {
	Id          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	FileName    string     `json:"file_name"`
	Extension   string     `json:"extension"`
	Xid         string     `json:"xid"`
	ContentType uint       `json:"content_type"`
}

type UserMediaDetail struct {
	Id           uint       `json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
	OwnerId      uint       `json:"owner_id"`
	DocumentType uint       `json:"document_type"`
}

func Upload(ctx *gin.Context) {
	var (
		media        Media
		documentType uint
		tok          *authentication.Token
		err          error
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

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.ParseError,
		})
		return
	}
	mType, err := mimetype.DetectFile(file.Filename)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.TypeError,
		})
		return
	}
	if !utils.IsValidFile(mType.String()) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.TypeError,
		})
		return
	}

	media.FileName = file.Filename
	media.Extension = filepath.Ext(file.Filename)
	media.Xid = xid.New().String()

	err = ctx.SaveUploadedFile(file, utils.FILE_UPLOAD_DIR+media.Xid+media.Extension)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
		return
	}

	defer openedFile.Close()
	err = utils.CreateThumb(media.Xid, media.Extension, openedFile)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
		return
	}

	_, err = database.InsertOne(media)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbInsertError,
		})
		return
	}
	if utils.IsValidImage(mType.String()) {
		documentType = 0
	}
	if utils.IsValidDocument(mType.String()) {
		documentType = 1
	}
	if utils.IsValidVideo(mType.String()) {

	}
	err = SetUserMediaDetail(documentType, tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbInsertError,
		})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, media)
}

func SetUserMediaDetail(documentType uint, userId uint) (err error) {
	var (
		userMediaDetail UserMediaDetail
	)

	userMediaDetail.OwnerId = userId
	userMediaDetail.DocumentType = documentType
	_, err = database.InsertOne(userMediaDetail)
	if err != nil {
		return err
	}
	return
}
