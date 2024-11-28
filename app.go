package main

import (
	"net/http"

	"github.com/xineman/go-server/db"
	"github.com/xineman/go-server/entities/tracks"

	"github.com/labstack/echo/v4"
)

func main() {
	db.Init()
	defer db.DbPool.Close()

	e := echo.New()
	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	tracks.TrackRouter(e.Group("/tracks"))

	e.Logger.Fatal(e.Start(":1323"))
}
