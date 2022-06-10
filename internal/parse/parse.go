package parse

import (
	"github.com/valyala/fasthttp"
	"strconv"
)

func IntGetParameter(paramName string, ctx *fasthttp.RequestCtx) (int, error) {
	param := string(ctx.QueryArgs().Peek(paramName))
	if param == "" {
		return -1, nil
	}
	parsedValue, err := strconv.Atoi(param)
	if err != nil {
		return -1, err
	}

	return parsedValue, nil
}

func BoolGetParameter(paramName string, ctx *fasthttp.RequestCtx) (bool, error) {
	param := string(ctx.QueryArgs().Peek(paramName))
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
	param := string(ctx.QueryArgs().Peek(paramName))
	return param
}

func IntSlugParameter(paramName string, ctx *fasthttp.RequestCtx) (int, error) {
	param := ctx.UserValue(paramName).(string)

	parsedValue, err := strconv.Atoi(param)

	if err != nil {
		return -1, err
	}

	return parsedValue, nil
}
