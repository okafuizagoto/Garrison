package middleware

import "github.com/kataras/iris/v12"

// AppOriginMiddleware extracts X-App-Origin header and stores it in context values.
func AppOriginMiddleware(ctx iris.Context) {
	origin := ctx.GetHeader("X-App-Origin")
	ctx.Values().Set("app_origin", origin)
	ctx.Next()
}
