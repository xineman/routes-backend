package main

import (
	"net/http"

	"github.com/xineman/go-server/db"
	"github.com/xineman/go-server/entities/tracks"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger" // echo-swagger middleware
	_ "github.com/xineman/go-server/docs"        // docs is generated by Swag CLI, you have to import it.
)

// @title Swagger Example API
// @version 2.0
func main() {
	db.Init()
	defer db.DbPool.Close()

	e := echo.New()
	e.Static("/static", "static")

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	tracksGroup := e.Group("/tracks")
	{
		tracksGroup.GET("", tracks.GetAll)
		tracksGroup.POST("", tracks.Create)
		tracksGroup.DELETE("", tracks.Delete)
		tracksGroup.POST("/bulk", tracks.CreateBulk)
	}

	e.Logger.Fatal(e.Start(":1323"))
}