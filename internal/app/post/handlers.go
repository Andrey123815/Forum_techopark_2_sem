package post

import (
	"db-forum/internal/app/models"
	"db-forum/internal/app/post/postRepo"
	"db-forum/internal/parse"
	"db-forum/internal/responseDelivery"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

type Handlers struct {
	PostRepo postRepo.PostRepository
}

func handleInternalServerError(err error, ctx *fasthttp.RequestCtx) bool {
	if err != nil {
		responseDelivery.SendInternalServerError(ctx)
		return true
	}

	return false
}

func (h *Handlers) GetPostDetails(ctx *fasthttp.RequestCtx) {
	stringifyID := ctx.UserValue("id").(string)
	postID, err := strconv.Atoi(stringifyID)
	if err != nil {
		responseDelivery.SendError(fasthttp.StatusNotFound, "", ctx)
		return
	}

	relatedParamsArr := parse.StringGetParameter("related", ctx)
	related := strings.Split(relatedParamsArr, ",")

	postDetails, err := h.PostRepo.GetPostDetails(postID, related)
	if err != nil {
		responseDelivery.SendError(fasthttp.StatusNotFound, "", ctx)
		return
	}

	responseDelivery.SendResponse(fasthttp.StatusOK, postDetails, ctx)
}

func (h *Handlers) ChangePostMessage(ctx *fasthttp.RequestCtx) {
	stringifyID := ctx.UserValue("id").(string)
	postID, err := strconv.Atoi(stringifyID)
	if err != nil {
		responseDelivery.SendError(fasthttp.StatusNotFound, "", ctx)
		return
	}

	var newPost models.Post
	err = json.Unmarshal(ctx.PostBody(), &newPost)
	if handleInternalServerError(err, ctx) == true {
		return
	}

	var related []string
	DetailPostInfo, err := h.PostRepo.GetPostDetails(postID, related)
	if err != nil {
		responseDelivery.SendError(fasthttp.StatusNotFound, "", ctx)
		return
	}

	oldPost := DetailPostInfo.Post
	if newPost.Message == "" || oldPost.Message == newPost.Message {
		responseDelivery.SendResponse(fasthttp.StatusOK, oldPost, ctx)
		return
	}

	updatedPost, err := h.PostRepo.ChangePostMessage(postID, newPost)
	if err != nil {
		responseDelivery.SendError(fasthttp.StatusNotFound, "", ctx)
		return
	}

	responseDelivery.SendResponse(fasthttp.StatusOK, updatedPost, ctx)
	return
}
