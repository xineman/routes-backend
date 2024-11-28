package tracks

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getTrackFileName(fileName string, trackName string) string {
	slug := strings.ReplaceAll(trackName, " ", "-")
	extension := filepath.Ext(fileName)
	if len(extension) > 0 {
		extension = extension[1:]
	}
	return fmt.Sprintf("%v-%v.%s", slug, time.Now().UnixMilli(), extension)
}

func createStaticFolderIfNotExist(subfolder string) error {
	if subfolder == "" {
		return errors.New("subfolder should be provided")
	}
	path := filepath.Join("static", subfolder)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	fmt.Println("Made sure directory exists:", path)
	return nil
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
