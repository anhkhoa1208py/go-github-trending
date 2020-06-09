package main

import (
	"backend-github-trending/db"
	"backend-github-trending/handler"
	"backend-github-trending/log"
	"backend-github-trending/repository/repo_imp"
	"backend-github-trending/router"
	"fmt"
	"github.com/labstack/echo/v4"
	"os"
)

func init()  {
	fmt.Println("init package main")
	os.Setenv("APP_NAME", "github")
	log.InitLogger(false)
}

func main() {
	fmt.Println("main function")
	sql := &db.Sql{
		Host: "localhost",
		Port: 5432,
		UserName: "postgres",
		Password: "Sak@1208",
		DbName: "github",
	}

	sql.Connect()
	defer sql.Close()

	e := echo.New()

	userHandler := handler.UserHandler{
		UserRepo: repo_imp.NewUserRep(sql),
	}

	api := router.API{
		Echo: e,
		UserHandler: userHandler,
	}

	api.SetupRputer()

	e.Logger.Fatal(e.Start(":3000"))
}
