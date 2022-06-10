package service

import (
	"db-forum/internal/app/forum/forumRepo"
	"db-forum/internal/app/service/serviceRepo"
	"db-forum/internal/app/thread/threadRepo"
	"db-forum/internal/app/user/userRepo"
	"db-forum/internal/responseDelivery"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
)

type Handlers struct {
	ServiceRepo serviceRepo.ServiceRepository
	ThreadRepo  threadRepo.ThreadRepository
	ForumRepo   forumRepo.ForumRepository
	UserRepo    userRepo.UserRepository
	InfoLog     *log.Logger
	ErrorLog    *log.Logger
}

func handleInternalServerError(err error, ctx *fasthttp.RequestCtx) bool {
	if err != nil {
		fmt.Println(err)
		responseDelivery.SendInternalServerError(ctx)
		return true
	}

	return false
}

func (h *Handlers) GetSystemStatus(ctx *fasthttp.RequestCtx) {
	systemStatus := h.ServiceRepo.GetSystemStatus()
	responseDelivery.SendResponse(fasthttp.StatusOK, systemStatus, ctx)
}

func (h *Handlers) ClearSystem(ctx *fasthttp.RequestCtx) {
	err := h.ServiceRepo.ClearSystem()
	if err != nil {
		handleInternalServerError(err, ctx)
	}

	responseDelivery.SendSuccessStatus(ctx)
}
