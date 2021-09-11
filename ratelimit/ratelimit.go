package ratelimit

import (
	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
	"github.com/kataras/iris/v12"
)

/*
go get -u github.com/didip/tollbooth/v6
LimitHandler is a middleware that performs
rate-limiting given a "limiter" configuration.
*/
func LimitHandler(lmt *limiter.Limiter) iris.Handler {
	return func(ctx iris.Context) {
		httpError := tollbooth.LimitByRequest(lmt, ctx.ResponseWriter(), ctx.Request())
		if httpError != nil {
			ctx.StatusCode(httpError.StatusCode)
			_, _ = ctx.WriteString(httpError.Message)
			ctx.StopExecution()
			return
		}
		ctx.Next()
	}
}
