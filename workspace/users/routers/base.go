package routers

import (
	"users/controllers"
	"users/middleware"

	"github.com/gold-gym/gymkit/middleware/auth"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

func Setup(app *iris.Application, db *gorm.DB) {
	jwtAuth := auth.New(auth.Config{})

	authController := &controllers.AuthController{DB: db}
	userController := &controllers.UserController{DB: db}

	// Public routes
	app.PartyFunc("/auth", func(r iris.Party) {
		r.Post("/login", authController.Login)
		r.Get("/authenticate", authController.Authenticate)
	})

	// Protected routes
	app.PartyFunc("/users", func(r iris.Party) {
		r.Use(middleware.AppOriginMiddleware)
		r.Use(jwtAuth)

		r.Get("/", userController.GetList)
		r.Get("/{id:uint64}", userController.GetDetail)
		r.Post("/", userController.Register)
		r.Put("/{id:uint64}", userController.Update)
		r.Delete("/{id:uint64}", userController.Delete)
	})

	// Health check
	app.Get("/ping", func(ctx iris.Context) {
		ctx.JSON(map[string]string{"status": "ok", "service": "users"})
	})
}
