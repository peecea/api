package post

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"peec/database"
	"peec/internal/authentication"
	"peec/internal/utils"
	"peec/internal/utils/errx"
	"peec/internal/utils/state"
	"strconv"
	"time"
)

type Post struct {
	Id        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	PosterId  uint       `json:"poster_id"`
	MediaXid  string     `json:"media_xid"`
}

type UserPost struct {
	Id        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	PostId    uint       `json:"post_id"`
	UserId    uint       `json:"userId"`
}

func CreatePost(ctx *gin.Context) {
	var (
		err      error
		tok      *authentication.Token
		post     Post
		userPost UserPost
	)
	err = ctx.ShouldBindJSON(&post)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		//
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

	postId, err := database.InsertOne(post)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		return
	}

	userPost.PostId = postId
	userPost.UserId = tok.UserId

	_, err = database.InsertOne(userPost)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		return
	}

	ctx.JSON(http.StatusOK, post)
	return
}

func DeletePost(ctx *gin.Context) {
	var (
		tok  *authentication.Token
		err  error
		post Post
	)
	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	postId, err := strconv.Atoi(ctx.Param("postId"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		return
	}

	if tok.UserId == state.ZERO {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	post, err = GetPost(postId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbGetError,
		})
		return
	}

	err = database.Delete(post)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbDeleteError,
		})
		return
	}

	ctx.JSON(http.StatusOK, true)
	return
}

func GetSinglePost(ctx *gin.Context) {
	var (
		tok  *authentication.Token
		err  error
		post Post
	)
	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	postId, err := strconv.Atoi(ctx.Param("postId"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.Lambda(err),
		})
		return
	}

	if tok.UserId == state.ZERO {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	post, err = GetPost(postId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.DbGetError,
		})
		return
	}
	ctx.JSON(http.StatusOK, post)
	return
}

func GetPosts(ctx *gin.Context) {
	var (
		posts []Post
		err   error
	)
	err = database.Select(posts, `SELECT * FROM post`)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}
	ctx.JSON(http.StatusOK, posts)
	return
}

func GetUserPosts(ctx *gin.Context) {
	var (
		posts []Post
		err   error
		tok   *authentication.Token
	)

	tok, err = authentication.GetTokenDataFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}

	err = database.Select(posts, `SELECT * FROM post where poster_id = ? `, tok.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: errx.UnAuthorizedError,
		})
		return
	}
	ctx.JSON(http.StatusOK, posts)
	return
}

/*

	UTILS

*/

func GetPost(postId int) (post Post, err error) {
	err = database.Get(&post, `SELECT * FROM post where id = ? `, postId)
	if err != nil {
		return post, err
	}
	return post, nil
}
