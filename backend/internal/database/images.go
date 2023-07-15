package database

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (dbCfg *DbConfig) saveImage(item ToplistItem, listID int) error {
	if len(item.Image) == 0 {
		// remove image for this rank if there were one previously
		err := removeImgByRank(item.Rank)
		if err != nil {
			return err
		}
		return nil
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	imagesDir := filepath.Join(currentDir, "internal", "database", "images")
	toplistDir := filepath.Join(imagesDir, fmt.Sprintf("%d", listID))
	

	err = os.MkdirAll(toplistDir, os.ModePerm)
	if err != nil {
		return err
	}

	mimeType := http.DetectContentType(item.Image)
	ext := getExtensionFromMimeType(mimeType)
	if ext == "" {
        return errors.New("unsupported img type")
    }

	imagePath := filepath.Join(toplistDir, fmt.Sprintf("%d%s", item.Rank, ext))

	file, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(item.Image))
	if err != nil {
		return err
	}

	return nil
}

func getExtensionFromMimeType(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	default:
		return ""
	}
}

func removeImgByRank(rank int) error {
	return nil
}

func (dbCfg *DbConfig) setImagePaths(items []ToplistItem, listID int) ([]ToplistItem, error) {
	imageDirPath := fmt.Sprintf("./internal/database/images/%d/", listID)
	if _, err := os.Stat(imageDirPath); err == nil {
		imgFiles, err := os.ReadDir(imageDirPath)
		if err != nil {
			return []ToplistItem{}, err
		}
	
		imageMap := make(map[int]string) // maps rank to image path

        for _, file := range imgFiles {
            fileBase := filepath.Base(file.Name())
            fileWithoutExt := strings.TrimSuffix(fileBase, filepath.Ext(fileBase))
            rank, err := strconv.Atoi(fileWithoutExt)
            if err == nil {
                imageMap[rank] = file.Name()
            }
        }

        for i, item := range items {
            if path, exists := imageMap[item.Rank]; exists {
                items[i].ImagePath = path
            }
        }
	} else {
		if !os.IsNotExist(err) {
			return []ToplistItem{}, err
		}
	}

	return items, nil
}