package util

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
)

func GetIrisApp() *iris.Application {
	app := iris.New()
	app.Use(recover.New())
	return app
}
