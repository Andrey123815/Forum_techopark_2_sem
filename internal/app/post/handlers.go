package post

import (
	"db-forum/internal/app/forum/forumRepo"
	"db-forum/internal/app/models"
	"db-forum/internal/app/post/postRepo"
	"db-forum/internal/app/thread/threadRepo"
	"db-forum/internal/app/user/userRepo"
	"db-forum/internal/parse"
	"db-forum/internal/responseDelivery"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"strings"
)

type Handlers struct {
	PostRepo   postRepo.PostRepository
	ThreadRepo threadRepo.ThreadRepository
	ForumRepo  forumRepo.ForumRepository
	UserRepo   userRepo.UserRepository
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
}

func handleInternalServerError(err error, ctx *fasthttp.RequestCtx) bool {
	if err != nil {
		fmt.Println(err)
		responseDelivery.SendInternalServerError(ctx)
		return true
	}

	return false
}

func (h *Handlers) GetPostDetails(ctx *fasthttp.RequestCtx) {
	postID, _ := parse.Int64SlugParameter("id", ctx)
	relatedOneString := string(ctx.QueryArgs().Peek("related"))
	related := strings.Split(relatedOneString, ",")

	postDetails, _ := h.PostRepo.GetPostDetails(postID, related)
	if postDetails["post"].(models.Post).Author == "" {
		responseDelivery.SendError(fasthttp.StatusNotFound, "Can't find forum", ctx)
		return
	}

	responseDelivery.SendResponse(fasthttp.StatusOK, postDetails, ctx)
}

func (h *Handlers) ChangePostMessage(ctx *fasthttp.RequestCtx) {
	postID, _ := parse.Int64SlugParameter("id", ctx)

	newPost := models.Post{}
	err := json.Unmarshal(ctx.PostBody(), &newPost)
	if handleInternalServerError(err, ctx) == true {
		return
	}

	var related []string
	oldPost, err := h.PostRepo.GetPostDetails(postID, related)

	if newPost.Message == "" || newPost.Message == oldPost["post"].(models.Post).Message {
		responseDelivery.SendResponse(fasthttp.StatusOK, oldPost["post"], ctx)
		return
	}

	updatedPost, err := h.PostRepo.ChangePostMessage(postID, newPost.Message)
	if err != nil {
		responseDelivery.SendError(fasthttp.StatusNotFound, "Can't find post", ctx)
		return
	}

	responseDelivery.SendResponse(fasthttp.StatusOK, updatedPost, ctx)
}
