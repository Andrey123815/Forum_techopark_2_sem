package service

import (
	"db-forum/internal/app/service/serviceRepo"
	"db-forum/internal/responseDelivery"
	"github.com/valyala/fasthttp"
)

type Handlers struct {
	ServiceRepo serviceRepo.ServiceRepository
}

func handleInternalServerError(err error, ctx *fasthttp.RequestCtx) bool {
	if err != nil {
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
