package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"users/config"
	"users/util"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var host, port string
	flag.StringVar(&host, "host", os.Getenv("HOST"), "host of the service")
	flag.StringVar(&port, "port", os.Getenv("PORT"), "port of the service")
	flag.Parse()

	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "8080"
	}

	db := util.GetDBConnection()
	defer func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
		log.Println("Closing database connection...")
	}()

	fmt.Printf("Service: %s\nVersion: %s\n", os.Getenv("APP_NAME"), os.Getenv("APP_VER"))

	routes := &config.Routes{DB: db}
	routes.Setup(host, port)
}
