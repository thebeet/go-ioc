package main

import (
	"database/sql"
	"example/component/db"
	"example/service/order"
	_ "example/service/order-impl"
	"example/service/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thebeet/go-ioc/pkg/ioc"
)

type App struct {
	UserService  user.Service  `autowire:""`
	OrderService order.Service `autowire:""`

	DbUser  *sql.DB `autowire:"db_user"`
	DbOrder *sql.DB `autowire:"db_order"`
}

func main() {
	var app App
	ioc.RegisterInstanceWithName(db.NewDb("user:userpass@tcp(127.0.0.1:8306)/user"), "db_user")
	ioc.RegisterInstanceWithName(db.NewDb("order:orderpass@tcp(127.0.0.1:8307)/order"), "db_order")
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

	r.GET("/order", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Query("id"))
		order := app.OrderService.GetOrderById(id)
		if order.Id == 0 {
			c.JSON(http.StatusNotFound, order)
		} else {
			c.JSON(http.StatusOK, order)
		}
	})

	r.GET("/health", func(c *gin.Context) {
		if app.DbUser.Ping() == nil && app.DbOrder.Ping() == nil {
			c.String(http.StatusOK, "ok")
		} else {
			c.String(http.StatusServiceUnavailable, "db fail")
		}
	})

	r.Run()
}
