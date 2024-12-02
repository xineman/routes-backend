package tracks

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5"
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

// @Summary		Save a new track
// @ID					create-track
// @Tags				Tracks
// @Accept			json
// @Param track formData file true "GPX track file"
// @Param photos formData []file true "Track photos"
// @Param name formData string true "Track name"
// @Produce		json
// @Success		201	{object}	Track
// @Failure		500	{object}	error
// @Router			/tracks [post]
func Create(c echo.Context) error {
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
}

// @Summary		Get all tracks
// @ID					get-all-tracks
// @Tags				Tracks
// @Accept			json
// @Produce		json
// @Success		200	{array}	Track
// @Failure		500	{object}	error
// @Router			/tracks [get]
func GetAll(c echo.Context) error {
	data, err := getTracks()
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusOK, []interface{}{})
		}
		return err
	}
	if len(data) != 0 {
		return c.JSON(http.StatusOK, data)
	}
	return c.JSON(http.StatusOK, []interface{}{})
}

// @Summary		Delete track
// @ID					delete-track
// @Tags				Tracks
// @Accept			json
// @Param id query string true "Track ID"
// @Produce		json
// @Success		200
// @Failure		404
// @Failure		500	{object}	error
// @Router			/tracks [delete]
func Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		fmt.Printf("Could not parse track id %v\n", id)
		return err
	}
	err = delete(id)
	if err != nil {
		fmt.Printf("Could not delete track %v\n", id)
		if err == pgx.ErrNoRows {
			return c.NoContent(http.StatusNotFound)
		}
	}
	return err
}

func CreateBulk(c echo.Context) error {
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
}
