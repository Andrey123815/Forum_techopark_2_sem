package responseDelivery

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type ErrorMessage struct {
	Message string
}

func SendSuccessStatus(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func SendInternalServerError(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
}

func SendError(errCode int, errorMsg string, ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(errCode)

	var errorMessage = ErrorMessage{errorMsg}

	JSONEncodedData, err := json.Marshal(errorMessage)
	if err != nil {
		SendInternalServerError(ctx)
	}

	ctx.SetBody(JSONEncodedData)
}

func SendResponse(status int, data interface{}, ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(status)

	JSONEncodedData, err := json.Marshal(data)
	if err != nil {
		SendInternalServerError(ctx)
	}

	ctx.SetBody(JSONEncodedData)
}
