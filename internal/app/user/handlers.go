package user

import (
	"db-forum/internal/app/models"
	"db-forum/internal/app/user/userRepo"
	"db-forum/internal/responseDelivery"
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type Handlers struct {
	UserRepo userRepo.UserRepository
}

func handleInternalServerError(err error, ctx *fasthttp.RequestCtx) bool {
	if err != nil {
		responseDelivery.SendInternalServerError(ctx)
		return true
	}

	return false
}

func handleUserError(statusCode int, err string, ctx *fasthttp.RequestCtx) {
	responseDelivery.SendError(statusCode, err, ctx)
}

func (h *Handlers) CreateNewUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	newUser := models.User{}
	err := json.Unmarshal(ctx.PostBody(), &newUser)
	if handleInternalServerError(err, ctx) == true {
		return
	}

	users, err := h.UserRepo.CheckUserExists(nickname, newUser.Email)
	if len(users) != 0 {
		responseDelivery.SendResponse(fasthttp.StatusConflict, users, ctx)
		return
	}

	userDetails, err := h.UserRepo.CreateNewUser(nickname, newUser.Fullname, newUser.About, newUser.Email)

	if handleInternalServerError(err, ctx) == true {
		return
	}

	responseDelivery.SendResponse(fasthttp.StatusCreated, userDetails, ctx)
}

func (h *Handlers) GetUserProfile(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	user, err := h.UserRepo.GetUserByNickname(nickname)

	if user == (models.User{}) || err != nil {
		handleUserError(fasthttp.StatusNotFound, "Can't find user by nickname: "+nickname, ctx)
		return
	}

	responseDelivery.SendResponse(fasthttp.StatusOK, user, ctx)
}

func (h *Handlers) UpdateUserProfile(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	userByNickname, _ := h.UserRepo.GetUserByNickname(nickname)
	if userByNickname.Nickname == "" {
		responseDelivery.SendError(fasthttp.StatusNotFound, "Can't find user by nickname: "+nickname, ctx)
		return
	}

	userNewSettings := models.User{}
	err := json.Unmarshal(ctx.PostBody(), &userNewSettings)
	if handleInternalServerError(err, ctx) == true {
		return
	}

	if userNewSettings.Nickname == "" && userNewSettings.About == "" && userNewSettings.Email == "" && userNewSettings.Fullname == "" {
		responseDelivery.SendResponse(fasthttp.StatusOK, userByNickname, ctx)
		return
	}

	if userNewSettings.Email != "" {
		userByEmail, err := h.UserRepo.GetUserByEmail(userNewSettings.Email)
		if err != nil || userByEmail.Email != "" {
			responseDelivery.SendError(fasthttp.StatusConflict, "This email is already registered by user: "+userByEmail.Nickname, ctx)
			return
		}
	}

	userWithNewSettings, err := h.UserRepo.UpdateUserProfile(nickname,
		userNewSettings.Fullname, userNewSettings.About, userNewSettings.Email)

	responseDelivery.SendResponse(fasthttp.StatusOK, userWithNewSettings, ctx)
}
