package tracks

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

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

		go saveFile(*track, fileName)
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
				errors <- saveFile(*file, "")
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

func getTrackFileName(fileName string, trackName string) string {
	slug := strings.ReplaceAll(trackName, " ", "-")
	extension := filepath.Ext(fileName)
	if len(extension) > 0 {
		extension = extension[1:]
	}
	return fmt.Sprintf("%v-%v.%s", slug, time.Now().UnixMilli(), extension)
}

func saveFile(file multipart.FileHeader, name string) error {
	fmt.Println("Processing file:", file.Filename)
	time.Sleep(time.Second * 2)
	src, err := file.Open()
	if err != nil {
		fmt.Println("Open error:", err)
		return err
	}
	defer src.Close()

	fileName := name
	if name == "" {
		fileName = fmt.Sprintf("%v-%v", time.Now().UnixMilli(), fileName)
	}
	dst, err := os.Create(filepath.Join("static/tracks", fileName))
	if err != nil {
		fmt.Println("Create error:", err, fileName)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		fmt.Println("Copy error:", err, fileName)
		return err
	}
	fmt.Println("Saved file:", fileName)
	return nil
}
