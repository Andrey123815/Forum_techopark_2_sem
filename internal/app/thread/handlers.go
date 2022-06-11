package thread

import (
	"db-forum/internal/app/models"
	"db-forum/internal/app/thread/threadRepo"
	"db-forum/internal/app/user/userRepo"
	"db-forum/internal/customErrors"
	"db-forum/internal/parse"
	"db-forum/internal/responseDelivery"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

type Handlers struct {
	ThreadRepo threadRepo.ThreadRepository
	UserRepo   userRepo.UserRepository
}

func handleInternalServerError(err error, ctx *fasthttp.RequestCtx) bool {
	if err != nil {
		responseDelivery.SendInternalServerError(ctx)
		return true
	}

	return false
}

func (h *Handlers) CreateNewPosts(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)

	searchedThread, err := h.ThreadRepo.GetThreadBySlugOrID(slugOrID)
	if searchedThread.Title == "" {
		responseDelivery.SendError(fasthttp.StatusNotFound, "Can't find post thread by id", ctx)
		return
	}

	newPosts := make([]models.Post, 0)
	err = json.Unmarshal(ctx.PostBody(), &newPosts)
	if handleInternalServerError(err, ctx) == true {
		return
	}

	if len(newPosts) == 0 {
		responseDelivery.SendResponse(fasthttp.StatusCreated, newPosts, ctx)
		return
	}

	createdPosts, err := h.ThreadRepo.CreateNewPosts(newPosts, searchedThread)
	if err != nil {
		switch err {
		case customErrors.NoAuthor:
			responseDelivery.SendError(fasthttp.StatusNotFound, err.Error(), ctx)
			return

		case customErrors.AnotherThreadPost:
			responseDelivery.SendError(fasthttp.StatusConflict, err.Error(), ctx)
			return
		}
	}

	responseDelivery.SendResponse(fasthttp.StatusCreated, createdPosts, ctx)
}

func (h *Handlers) VoteThread(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)

	thread, err := h.ThreadRepo.GetThreadBySlugOrID(slugOrID)
	if err != nil {
		responseDelivery.SendError(fasthttp.StatusNotFound, "", ctx)
		return
	}

	var vote models.Vote
	err = json.Unmarshal(ctx.PostBody(), &vote)
	if handleInternalServerError(err, ctx) == true {
		return
	}

	_, err = h.UserRepo.GetUserByNickname(vote.Nickname)
	if err != nil {
		responseDelivery.SendError(fasthttp.StatusNotFound, "user not found", ctx)
		return
	}

	err = h.ThreadRepo.VoteThread(vote.Nickname, thread.Id, vote.Voice, &thread)

	responseDelivery.SendResponse(fasthttp.StatusOK, thread, ctx)
}

func (h *Handlers) GetThreadDetails(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)

	threadDetails, err := h.ThreadRepo.GetThreadBySlugOrID(slugOrID)
	if err != nil {
		responseDelivery.SendError(fasthttp.StatusNotFound, "Can't find thread by slug", ctx)
		return
	}

	responseDelivery.SendResponse(fasthttp.StatusOK, threadDetails, ctx)
}

func (h *Handlers) GetThreadPosts(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)

	thread, err := h.ThreadRepo.GetThreadBySlugOrID(slugOrID)
	if err != nil {
		responseDelivery.SendError(fasthttp.StatusNotFound, fmt.Sprintf("Can't find user with slug #%s\n", slugOrID), ctx)
		return
	}

	limit := parse.StringGetParameter("limit", ctx)
	since := parse.StringGetParameter("since", ctx)
	sortType := parse.StringGetParameter("sort", ctx)
	desc := parse.StringGetParameter("desc", ctx)

	sortDirection := "ASC"
	if desc == "true" {
		sortDirection = "DESC"
	}

	if sortType == "" {
		sortType = "flat"
	}

	posts, err := h.ThreadRepo.GetThreadPosts(thread.Id, limit, since, sortType, sortDirection)
	if err != nil {
		responseDelivery.SendError(fasthttp.StatusNotFound, "posts not found", ctx)
		return
	}

	responseDelivery.SendResponse(fasthttp.StatusOK, posts, ctx)
}

func (h *Handlers) UpdateThreadDetails(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)

	thread, err := h.ThreadRepo.GetThreadBySlugOrID(slugOrID)
	if err != nil {
		responseDelivery.SendError(fasthttp.StatusNotFound, "can't find user", ctx)
		return
	}

	var threadDetails models.Thread
	err = json.Unmarshal(ctx.PostBody(), &threadDetails)
	if handleInternalServerError(err, ctx) == true {
		return
	}

	if threadDetails.Title == "" && threadDetails.Message == "" {
		responseDelivery.SendResponse(fasthttp.StatusOK, thread, ctx)
		return
	}

	if threadDetails.Title != "" {
		thread.Title = threadDetails.Title
	}
	if threadDetails.Message != "" {
		thread.Message = threadDetails.Message
	}

	updatedThread, err := h.ThreadRepo.UpdateThreadDetails(thread.Id, thread.Title, thread.Message)
	if updatedThread == (models.Thread{}) {
		responseDelivery.SendResponse(fasthttp.StatusNotFound, "can't find user", ctx)
		return
	}

	responseDelivery.SendResponse(fasthttp.StatusOK, updatedThread, ctx)
}
