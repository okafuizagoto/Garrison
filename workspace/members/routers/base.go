package routers

import (
	"members/controllers"
	"members/middleware"

	"github.com/gold-gym/gymkit/middleware/auth"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

func Setup(app *iris.Application, db *gorm.DB) {
	jwtAuth := auth.New(auth.Config{})
	memberController := &controllers.MemberController{DB: db}

	app.PartyFunc("/members", func(r iris.Party) {
		r.Use(middleware.AppOriginMiddleware)
		r.Use(jwtAuth)

		r.Get("/", memberController.GetList)
		r.Get("/{id:uint64}", memberController.GetDetail)
		r.Post("/", memberController.Create)
		r.Put("/{id:uint64}", memberController.Update)
		r.Delete("/{id:uint64}", memberController.Delete)
	})

	app.Get("/ping", func(ctx iris.Context) {
		ctx.JSON(map[string]string{"status": "ok", "service": "members"})
	})
}
