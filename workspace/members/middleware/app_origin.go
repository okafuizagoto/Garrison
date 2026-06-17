package middleware

import "github.com/kataras/iris/v12"

func AppOriginMiddleware(ctx iris.Context) {
	origin := ctx.GetHeader("X-App-Origin")
	ctx.Values().Set("app_origin", origin)
	ctx.Next()
}
