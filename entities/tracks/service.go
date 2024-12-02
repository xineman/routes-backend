package tracks

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	ph "github.com/xineman/go-server/entities/photos/repo"
	tr "github.com/xineman/go-server/entities/tracks/repo"
	"github.com/xineman/go-server/utils"
)

type PhotoResult struct {
	err      error
	fileName string
}

func processTrack(createTrackDto CreateTrack) (tr.Track, error) {
	fileName := getTrackFileName(createTrackDto.Track.Filename, createTrackDto.Name)

	newTrack, err := tr.SaveTrack(tr.TrackMetadata{createTrackDto.Name, fileName})
	if err != nil {
		fmt.Println("saveTrack error:", err)
		return tr.Track{}, err
	}

	go saveTrackFiles(createTrackDto, newTrack)
	return newTrack, nil
}

func saveTrackFiles(createTrackDto CreateTrack, track tr.Track) {
	err := utils.CreateStaticFolderIfNotExist("tracks")
	if err != nil {
		fmt.Println("Could not create static folder", err)
		return
	}

	err = utils.SaveFile(*createTrackDto.Track, fmt.Sprintf("tracks/%s", track.FileName))
	if err != nil {
		fmt.Println("Could not save file", err)
		return
	}

	err = utils.CreateStaticFolderIfNotExist("images")
	if err != nil {
		fmt.Println("Could not create static folder", err)
		return
	}

	numberOfPhotos := len(createTrackDto.Photos)
	var photoResults = make(chan PhotoResult, numberOfPhotos)
	var wg sync.WaitGroup
	wg.Add(numberOfPhotos)
	for _, photo := range createTrackDto.Photos {
		go func() {
			defer wg.Done()
			fileName := getPhotoFileName(photo.Filename)
			err := utils.SaveFile(*photo, fmt.Sprintf("images/%s", fileName))
			photoResults <- PhotoResult{err, fileName}
		}()
	}
	go func() {
		wg.Wait()
		close(photoResults)
	}()

	failedPhotos := []string{}
	successfulPhotos := []string{}
	for photoResult := range photoResults {
		if photoResult.err != nil {
			failedPhotos = append(failedPhotos, fmt.Sprintf("%s: %s", photoResult.fileName, photoResult.err))
		} else {
			successfulPhotos = append(successfulPhotos, photoResult.fileName)
		}
	}

	ph.SavePhotos(track.Id, successfulPhotos)
	if len(failedPhotos) > 0 {
		fmt.Println("Some photos could not be saved:")
		for _, message := range failedPhotos {
			fmt.Println(message)
		}
	}
}

func getTrackFileName(fileName string, trackName string) string {
	slug := strings.ReplaceAll(trackName, " ", "-")
	extension := filepath.Ext(fileName)
	if len(extension) > 0 {
		extension = extension[1:]
	}
	return fmt.Sprintf("%v-%v.%s", slug, time.Now().UnixMilli(), extension)
}

func getPhotoFileName(fileName string) string {
	extension := filepath.Ext(fileName)
	base := filepath.Base(fileName)

	slug := strings.ReplaceAll(strings.TrimSuffix(base, extension), " ", "-")
	if len(extension) > 0 {
		extension = extension[1:]
	}
	return fmt.Sprintf("%v-%v.%s", slug, time.Now().UnixMilli(), extension)
}

func delete(id int) error {
	track, err := tr.GetTrack(id)
	if err != nil {
		return err
	}
	go utils.DeleteFile(filepath.Join("static/tracks", track.FileName))
	for _, ph := range track.Photos {
		go utils.DeleteFile(filepath.Join("static/images", ph))
	}
	return tr.DeleteTrack(id)
}
