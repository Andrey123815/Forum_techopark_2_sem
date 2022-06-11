package forum

import (
	"db-forum/internal/app/forum/forumRepo"
	"db-forum/internal/app/models"
	"db-forum/internal/app/user/userRepo"
	"db-forum/internal/parse"
	"db-forum/internal/responseDelivery"
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type Handlers struct {
	ForumRepo forumRepo.ForumRepository
	UserRepo  userRepo.UserRepository
}

func handleInternalServerError(err error, ctx *fasthttp.RequestCtx) bool {
	if err != nil {
		responseDelivery.SendInternalServerError(ctx)
		return true
	}

	return false
}

func (h *Handlers) CreateNewForum(ctx *fasthttp.RequestCtx) {
	newForum := models.Forum{}
	err := json.Unmarshal(ctx.PostBody(), &newForum)
	if handleInternalServerError(err, ctx) == true {
		return
	}

	forumFound, err := h.ForumRepo.GetForumBySlug(newForum.Slug)
	if forumFound.Slug != "" {
		responseDelivery.SendResponse(fasthttp.StatusConflict, forumFound, ctx)
		return
	}

	userByNickname, err := h.UserRepo.GetUserByNickname(newForum.User)
	if userByNickname == (models.User{}) {
		responseDelivery.SendError(fasthttp.StatusNotFound, "Can't find user with nickname: "+newForum.User, ctx)
		return
	}
	forumDetails, err := h.ForumRepo.CreateNewForum(newForum.Title, userByNickname.Nickname, newForum.Slug)

	responseDelivery.SendResponse(fasthttp.StatusCreated, forumDetails, ctx)
}

func (h *Handlers) GetForum(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	forumDetails, err := h.ForumRepo.GetForumBySlug(slug)
	if err != nil {
		responseDelivery.SendError(fasthttp.StatusNotFound, "Can't find forum with slug: "+slug, ctx)
		return
	}

	responseDelivery.SendResponse(fasthttp.StatusOK, forumDetails, ctx)
}

func (h *Handlers) CreateNewThread(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	forumBySlug, err := h.ForumRepo.GetForumBySlug(slug)
	if err != nil {
		responseDelivery.SendError(fasthttp.StatusNotFound, "Can't find forum by slug: "+slug, ctx)
		return
	}

	newThread := models.Thread{}
	err = json.Unmarshal(ctx.PostBody(), &newThread)
	if handleInternalServerError(err, ctx) == true {
		return
	}

	newThread.Forum = forumBySlug.Slug

	existThread, err := h.ForumRepo.GetAlreadyExistThread(newThread.Slug)
	if existThread != (models.Thread{}) {
		responseDelivery.SendResponse(fasthttp.StatusConflict, existThread, ctx)
		return
	}

	author, err := h.UserRepo.GetUserByNickname(newThread.Author)
	if author == (models.User{}) {
		responseDelivery.SendError(fasthttp.StatusNotFound, "Can't find thread author by nickname: "+newThread.Author, ctx)
		return
	}

	createdThread, err := h.ForumRepo.CreateNewThread(newThread.Title, newThread.Author,
		newThread.Forum, newThread.Message, newThread.Slug, newThread.Created)

	responseDelivery.SendResponse(fasthttp.StatusCreated, createdThread, ctx)
}

func (h *Handlers) GetThreads(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	_, err := h.ForumRepo.GetForumBySlug(slug)
	if err != nil {
		responseDelivery.SendError(fasthttp.StatusNotFound, "Can't find forum with slug: "+slug, ctx)
		return
	}

	limit := parse.StringGetParameter("limit", ctx)
	desc := parse.StringGetParameter("desc", ctx)
	since := parse.StringGetParameter("since", ctx)

	threads, err := h.ForumRepo.GetThreadsBySlug(slug, limit, desc, since)

	responseDelivery.SendResponse(fasthttp.StatusOK, threads, ctx)
}

func (h *Handlers) GetForumUsers(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	forum, err := h.ForumRepo.GetForumBySlug(slug)
	if err != nil {
		responseDelivery.SendError(fasthttp.StatusNotFound, "Can't find forum with slug: "+slug, ctx)
		return
	}

	limit := parse.StringGetParameter("limit", ctx)
	desc := parse.StringGetParameter("desc", ctx)
	since := parse.StringGetParameter("since", ctx)

	sortDirection := "ASC"
	if desc == "true" {
		sortDirection = "DESC"
	}

	users, err := h.ForumRepo.GetForumUsers(forum.Id, limit, sortDirection, since)

	responseDelivery.SendResponse(fasthttp.StatusOK, users, ctx)
}
