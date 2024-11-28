package tracks

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/xineman/go-server/utils"
)

type (
	Response struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
	}

	CreateTrack struct {
		Name   string `form:"name"`
		Track  *multipart.FileHeader
		Photos []*multipart.FileHeader
	}
)

func TrackRouter(group *echo.Group) {
	group.GET("", func(c echo.Context) error {
		data, err := getTracks()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, data)
	})

	group.POST("", func(c echo.Context) error {
		var createTrackDto CreateTrack

		if err := c.Bind(&createTrackDto); err != nil {
			fmt.Println("Could not bind data", err)
			return err
		}
		form, err := c.MultipartForm()
		if err != nil {
			fmt.Println("MultipartForm error:", err)
			return err
		}
		createTrackDto.Track = form.File["track"][0]
		createTrackDto.Photos = form.File["photos"]

		newTrack, err := processTrack(createTrackDto)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, newTrack)
	})

	group.POST("/bulk", func(c echo.Context) error {
		form, err := c.MultipartForm()
		if err != nil {
			fmt.Println("MultipartForm error:", err)
			return err
		}

		var wg sync.WaitGroup
		wg.Add(len(form.File["tracks"]))

		var errors = make(chan error)

		for _, file := range form.File["tracks"] {
			fmt.Println("Start job")
			go func() {
				defer wg.Done()
				err := utils.CreateStaticFolderIfNotExist("tracks")
				if err != nil {
					errors <- err
				} else {
					errors <- utils.SaveFile(*file, fmt.Sprintf("tracks/%s", file.Filename))
				}
			}()
		}

		go func() {
			wg.Wait()
			close(errors)
		}()

		errorsCount := 0
		for err := range errors {
			if err != nil {
				errorsCount++
			}
		}

		var message = "All files were saved"
		if errorsCount > 0 {
			message = fmt.Sprintf("Some files could not be saved, errors: %d", errorsCount)
		}

		return c.JSON(http.StatusOK, Response{Success: true, Message: message})
	})
}
