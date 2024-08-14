package video

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/joinverse/xid"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"peec/database"
	"peec/internal/authentication"
	"peec/internal/configuration"
	"peec/internal/utils"
	"peec/internal/utils/errx"
	"peec/internal/utils/state"
	"peec/pkg/media"
	"time"
)

const (
	Presentation = utils.VideoPresentation
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
	DocumentXid  string     `json:"document_xid"`
}

func UploadVideo(ctx *gin.Context) {
	var (
		media Media
		tok   *authentication.Token
		err   error
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

	mType, err := DetectMimeType(file)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.TypeError,
		})
		return
	}

	if !utils.IsValidVideo(mType.String()) {
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

	_, err = database.InsertOne(media)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbInsertError,
		})
		return
	}

	err = SetUserMediaDetail(tok.UserId, media.Xid)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbInsertError,
		})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, media)
}

func GetProfileVideo(ctx *gin.Context) {
	var (
		err   error
		media Media
		tok   *authentication.Token
	)

	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
		return
	}

	err = database.Get(&media, `SELECT media.*
FROM media
         JOIN user_media_detail ON media.xid = user_media_detail.document_xid
         JOIN user ON user.id = user_media_detail.owner_id
WHERE user_media_detail.owner_id = ? AND user_media_detail.document_type = ?`, tok.UserId, Presentation)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbGetError,
		})
		return
	}

	networkLink := "http://" + configuration.App.Host + ":" + configuration.App.Port + "/api/public/" + media.Xid + media.Extension

	ctx.JSON(http.StatusOK, networkLink)
	return
}

func UpdateProfileVideo(ctx *gin.Context) {
	var (
		media Media
		tok   *authentication.Token
		err   error
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
			Message: err,
		})
		return
	}

	oldMedia, err := GetCurrentUserVideo(tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
		return
	}

	media.FileName = file.Filename
	media.Extension = filepath.Ext(file.Filename)
	media.Xid = oldMedia.Xid

	err = RemoveCurrentUserVideo(oldMedia)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
	}

	_, err = database.InsertOne(media)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbInsertError,
		})
	}
	ctx.AbortWithStatus(http.StatusOK)
}

func RemoveProfileVideo(ctx *gin.Context) {
	var (
		media media.Media
		tok   *authentication.Token
		err   error
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

	media, err = GetCurrentUserVideo(tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	err = RemoveCurrentUserVideo(media)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbDeleteError,
		})
		return
	}

	userMediaDetail, err := GetUserMediaDetail(tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbDeleteError,
		})
		return
	}

	err = RemoveUserMediaDetail(userMediaDetail)
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

func DetectMimeType(file *multipart.FileHeader) (mType *mimetype.MIME, err error) {
	readFile, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	mType, err = mimetype.DetectReader(readFile)
	if err != nil {
		return nil, err
	}
	return mType, nil
}

func SetUserMediaDetail(userId uint, documentXid string) (err error) {
	var (
		userMediaDetail UserMediaDetail
	)
	userMediaDetail.OwnerId = userId
	userMediaDetail.DocumentType = Presentation
	userMediaDetail.DocumentXid = documentXid
	_, err = database.InsertOne(userMediaDetail)
	if err != nil {
		return err
	}
	return
}

func GetUserMediaDetail(userId uint) (userMediaDetail UserMediaDetail, err error) {
	err = database.Get(&userMediaDetail, `SELECT user_media_detail.* FROM  user_media_detail WHERE user_media_detail.owner_id =? `, userId)
	if err != nil {
		return userMediaDetail, err
	}
	return userMediaDetail, err
}

func GetCurrentUserVideo(userId uint) (media media.Media, err error) {
	err = database.Get(&media, `SELECT media.*
FROM media
         JOIN user_media_detail ON media.xid = user_media_detail.document_xid
         JOIN user ON user.id = user_media_detail.owner_id
WHERE user_media_detail.owner_id = ? AND user_media_detail.document_type = ?`, userId, Presentation)
	if err != nil {
		return media, err
	}
	return media, err
}

func RemoveUserMediaDetail(userMediaDetail UserMediaDetail) (err error) {
	err = database.Delete(userMediaDetail)
	if err != nil {
		return err
	}
	return
}

func RemoveCurrentUserVideo(media media.Media) (err error) {
	err = database.Delete(media)
	if err != nil {
		return err
	}
	return
}
