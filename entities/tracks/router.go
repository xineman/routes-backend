package tracks

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func TrackRouter(group *echo.Group) {
	group.GET("", func(c echo.Context) error {
		data, err := getTracks()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, data)
	})

	group.POST("", func(c echo.Context) error {
		track, err := c.FormFile("track")
		if err != nil {
			fmt.Println("FormFile error:", err)
			return err
		}
		name := c.FormValue("name")
		fileName := getTrackFileName(track.Filename, name)

		go func() {
			err := createStaticFolderIfNotExist("tracks")
			if err != nil {
				fmt.Println("Could not create static folder", err)
			} else {
				saveFile(*track, fileName)
			}
		}()
		newTrack, err := saveTrack(TrackMetadata{name, fileName})
		if err != nil {
			fmt.Println("saveTrack error:", err)
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
				err := createStaticFolderIfNotExist("tracks")
				if err != nil {
					errors <- err
				} else {
					errors <- saveFile(*file, "")
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
