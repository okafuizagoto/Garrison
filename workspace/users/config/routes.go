package config

import (
	"fmt"
	"log"
	"os"
	"users/routers"
	"users/util"

	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

type Routes struct {
	DB *gorm.DB
}

func (r *Routes) Setup(host, port string) {
	app := util.GetIrisApp()

	routers.Setup(app, r.DB)

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("users service starting on %s", addr)

	if err := app.Listen(addr, iris.WithoutServerError(iris.ErrServerClosed)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
