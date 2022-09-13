package main

import (
	"database/sql"
	"example/db"
	"example/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thebeet/go-ioc/pkg/ioc"
)

type App struct {
	UserService user.Service `autowire:""`
	DbUser      *sql.DB      `autowire:"db_user"`
}

func main() {
	var app App
	ioc.RegisterInstance(db.NewDb("user:userpass@tcp(127.0.0.1:8306)/user"))
	ioc.Fill(&app)

	r := gin.Default()
	r.GET("/user", func(c *gin.Context) {
		user := app.UserService.GetUserByName(c.Query("name"))
		if user.Id == 0 {
			c.JSON(http.StatusNotFound, user)
		} else {
			c.JSON(http.StatusOK, user)
		}
	})

	r.GET("/health", func(c *gin.Context) {
		if app.DbUser.Ping() == nil {
			c.String(http.StatusOK, "ok")
		} else {
			c.String(http.StatusServiceUnavailable, "db fail")
		}
	})

	r.Run()
}
