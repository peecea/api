package profile

import (
	"github.com/gabriel-vasile/mimetype"
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
	"peec/pkg/user"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	UserProfileImage = utils.UserProfileImage
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

func Upload(ctx *gin.Context) {
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

	if !utils.IsValidImage(mType.String()) {
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

	defer func(openedFile multipart.File) {
		err := openedFile.Close()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
				Message: err,
			})
			return
		}
	}(openedFile)

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

	err = UpdateUserProfileImageXid(tok.UserId, media.Xid)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
		return
	}

	err = SetUserMediaDetail(tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbInsertError,
		})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, media)
}

func GetProfileImage(ctx *gin.Context) {
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
         JOIN user ON user.profile_image_xid = media.xid
         JOIN user_media_detail ON user.id = user_media_detail.owner_id
WHERE user_media_detail.owner_id = ? AND document_type = ?`, tok.UserId, UserProfileImage)
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

func GetProfileThumb(ctx *gin.Context) {
	var (
		err        error
		mediaThumb utils.MediaThumb
		tok        *authentication.Token
	)

	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	mediaThumb, err = GetCurrentUserProfileThumb(tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbGetError,
		})
		return
	}

	networkLink := "http://" + configuration.App.Host + ":" + configuration.App.Port + "/api/public/thumb/" + mediaThumb.MediaXid + mediaThumb.Extension

	ctx.JSON(http.StatusOK, networkLink)
	return
}

func UpdateProfileImage(ctx *gin.Context) {
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

	oldMedia, err := GetCurrentUserProfile(tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
		return
	}

	media.FileName = file.Filename
	media.Extension = filepath.Ext(file.Filename)
	media.Xid = oldMedia.Xid

	err = RemoveCurrentUserProfile(oldMedia)
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

func RemoveProfileImage(ctx *gin.Context) {
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
	media, err = GetCurrentUserProfile(tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	err = UpdateUserProfileImageXid(tok.UserId, "")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err,
		})
		return
	}

	err = RemoveCurrentUserProfile(media)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbDeleteError,
		})
		return
	}
	err = RemoveUserMediaDetail(tok.UserId)
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

func UpdateUserProfileImageXid(userId uint, xid string) (err error) {
	var (
		usr user.User
	)
	usr, err = GetCurrentUser(userId)
	usr.ProfileImageXid = xid
	err = database.Update(usr)
	if err != nil {
		return err
	}
	return
}

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

func SetUserMediaDetail(userId uint) (err error) {
	var (
		userMediaDetail utils.UserMediaDetail
	)

	userMediaDetail.OwnerId = userId
	userMediaDetail.DocumentType = UserProfileImage
	_, err = database.InsertOne(userMediaDetail)
	if err != nil {
		return err
	}
	return
}

func GetCurrentUserProfile(userId uint) (media media.Media, err error) {
	err = database.Get(&media, `SELECT media.*
FROM media
         JOIN user ON user.profile_image_xid = media.xid
         JOIN user_media_detail ON user.id = user_media_detail.owner_id
WHERE user_media_detail.owner_id = ? AND document_type = ?`, userId, UserProfileImage)
	if err != nil {
		return media, err
	}
	return media, err
}

func GetCurrentUser(userId uint) (user user.User, err error) {
	err = database.Get(&user, `SELECT * FROM user WHERE user.id = ?`, userId)
	if err != nil {
		return user, err
	}
	return user, err
}

func GetCurrentUserProfileThumb(userId uint) (media utils.MediaThumb, err error) {
	err = database.Get(&media, `SELECT media_thumb.*
				FROM media_thumb
						 JOIN media ON  media.xid = media_thumb.media_xid
						 JOIN user ON user.profile_image_xid = media.xid
						 JOIN user_media_detail ON user.id = user_media_detail.owner_id
				WHERE user_media_detail.owner_id = ? and document_type = ?`, userId, UserProfileImage)
	if err != nil {
		return media, err
	}
	return media, err
}

func GetUserMediaDetail(userId uint) (userMediaDetail utils.UserMediaDetail, err error) {
	err = database.Get(&userMediaDetail, `SELECT user_media_detail.* FROM  user_media_detail WHERE user_media_detail.owner_id =? `, userId)
	if err != nil {
		return userMediaDetail, err
	}
	return userMediaDetail, err
}

func RemoveUserMediaDetail(userId uint) (err error) {
	userMediaDetail, err := GetUserMediaDetail(userId)
	if err != nil {
		return err
	}
	err = database.Delete(userMediaDetail)
	if err != nil {
		return err
	}
	return
}

func RemoveCurrentUserProfile(media media.Media) (err error) {
	err = database.Delete(media)
	if err != nil {
		return err
	}
	return
}
