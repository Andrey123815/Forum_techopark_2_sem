package parse

import (
	"github.com/valyala/fasthttp"
	"strconv"
)

func BoolGetParameter(paramName string, ctx *fasthttp.RequestCtx) (bool, error) {
	param := string(ctx.FormValue(paramName))
	if param == "" {
		return false, nil
	}
	parsedValue, err := strconv.ParseBool(param)
	if err != nil {
		return false, err
	}

	return parsedValue, nil
}

func StringGetParameter(paramName string, ctx *fasthttp.RequestCtx) string {
	param := string(ctx.FormValue(paramName))
	return param
}

func Int64SlugParameter(paramName string, ctx *fasthttp.RequestCtx) (int64, error) {
	param := ctx.UserValue(paramName).(string)

	parsedValue, err := strconv.ParseInt(param, 10, 64)

	if err != nil {
		return -1, err
	}

	return parsedValue, nil
}
