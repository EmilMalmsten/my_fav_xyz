package database

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func getExtensionFromMimeType(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	default:
		return ""
	}
}

func saveImage(item *ToplistItem, listID int) error {

	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	exeDir := filepath.Dir(exePath)

	imagesDir := filepath.Join(exeDir, "images")
	toplistDir := filepath.Join(imagesDir, fmt.Sprintf("%d", listID))

	err = os.MkdirAll(toplistDir, os.ModePerm)
	if err != nil {
		return err
	}

	imagePath := item.ImagePath
	if imagePath == "" {
		fileName := uuid.New().String()
		mimeType := http.DetectContentType(item.Image)
		ext := getExtensionFromMimeType(mimeType)
		if ext == "" {
			return errors.New("unsupported img type")
		}
		imagePath = fmt.Sprintf("/%d/%s%s", listID, fileName, ext)
	}

	f, err := os.Create(filepath.Clean(imagesDir + imagePath))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, bytes.NewReader(item.Image))
	if err != nil {
		return err
	}
	item.ImagePath = imagePath
	return nil
}

func (dbCfg *DbConfig) deleteImage(item *ToplistItem, listID int) error {
	if item.ImagePath != "" {
		return nil
	}

	var imgPath string
	query := "SELECT image_path FROM list_items WHERE toplist_id = $1 AND rank = $2"

	row := dbCfg.database.QueryRowContext(context.Background(), query, listID, item.Rank)
	err := row.Scan(&imgPath)
	if err != nil {
		return err
	}
	if imgPath == "" {
		return nil
	}

	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	exeDir := filepath.Dir(exePath)

	imagesDir := filepath.Join(exeDir, "images")

	err = os.Remove(imagesDir + imgPath)
	if err != nil {
		return err
	}

	item.ImagePath = ""
	return nil
}

func deleteToplistImages(listID int) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	exeDir := filepath.Dir(exePath)

	imagesDir := filepath.Join(exeDir, "images")
	toplistDir := filepath.Join(imagesDir, fmt.Sprintf("%d", listID))

	err = os.RemoveAll(toplistDir)
	if err != nil {
		return err
	}

	return nil
}

func (dbCfg *DbConfig) saveOrDeleteImage(item *ToplistItem, listID int) error {

	if len(item.Image) > 0 {
		err := saveImage(item, listID)
		if err != nil {
			log.Println(err)
		}
	} else {
		err := dbCfg.deleteImage(item, listID)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
