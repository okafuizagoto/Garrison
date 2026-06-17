package utils

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
)

// NewApp creates a preconfigured Iris application.
func NewApp() *iris.Application {
	app := iris.New()
	app.Use(recover.New())
	return app
}
