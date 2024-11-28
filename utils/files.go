package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func CreateStaticFolderIfNotExist(subfolder string) error {
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

func SaveFile(file multipart.FileHeader, fullName string) error {
	fmt.Println("Processing file:", file.Filename)
	time.Sleep(time.Second * 2)
	src, err := file.Open()
	if err != nil {
		fmt.Println("Open error:", err)
		return err
	}
	defer src.Close()

	dst, err := os.Create(filepath.Join("static", fullName))
	if err != nil {
		fmt.Println("Create error:", err, fullName)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		fmt.Println("Copy error:", err, fullName)
		return err
	}
	fmt.Println("Saved file:", fullName)
	return nil
}
