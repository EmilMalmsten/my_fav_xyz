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

func getToplistDir(listID int) (string, error) {

	exePath, err := os.Executable()
    if err != nil {
		return "", err
    }

    exeDir := filepath.Dir(exePath)

	imagesDir := filepath.Join(exeDir, "images")
	toplistDir := filepath.Join(imagesDir, fmt.Sprintf("%d", listID))
	return toplistDir, nil
}

func saveImage(item ToplistItem, listID int) error {

	toplistDir, err := getToplistDir(listID)
	if err != nil {
		return err
	}

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

/*
func (dbCfg *DbConfig) deleteToplistImages(listID int) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	imagesDir := filepath.Join(currentDir, "internal", "database", "images")
	toplistDir := filepath.Join(imagesDir, fmt.Sprintf("%d", listID))
	filePaths, err := filepath.Glob(filepath.Join(toplistDir, "*"))
	if err != nil {
		return err
	}

	for _, filePath := range filePaths {
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	return nil
}
*/

func updateImage(item, swapItem *ToplistItem, listID int) error {

	toplistDir, err := getToplistDir(listID)
	if err != nil {
		return err
	}

	swapItemFullPath := filepath.Join(toplistDir, swapItem.ImagePath)
	currItemFullPath := filepath.Join(toplistDir, item.ImagePath)
	tempPath := filepath.Join(toplistDir, "temp.png")

	fmt.Printf("swapping %s to temp\n", swapItemFullPath)
	err = os.Rename(swapItemFullPath, tempPath)
	if err != nil {
		return err
	}

	fmt.Printf("swapping %s to %s\n", currItemFullPath, swapItemFullPath)
	err = os.Rename(currItemFullPath, swapItemFullPath)
	if err != nil {
		return err
	}

	fmt.Printf("swapping temp to %s\n", currItemFullPath)
	err = os.Rename(tempPath, currItemFullPath)
	if err != nil {
		return err
	}

	tmp := item.ImagePath
	item.ImagePath = swapItem.ImagePath
	swapItem.ImagePath = tmp

	return nil
}

func findSwapItem(items []ToplistItem, currentItem *ToplistItem) *ToplistItem {
	for i := range items {
		fileRank := getRankFromFileName(items[i].ImagePath)

		if currentItem.Rank == fileRank {
			// we found the item to swap with
			return &items[i]
		}
	}

	return nil
}

func handleImageChanges(items []ToplistItem, listID int) ([]ToplistItem, error) {
	for i := range items {

		if items[i].ImagePath == "" && len(items[i].Image) > 0 {
			err := saveImage(items[i], listID)
			if err != nil {
				fmt.Println(err)
				return []ToplistItem{}, err
			}
		} else if items[i].ImagePath != "" && len(items[i].Image) == 0 {
			//  check if item needs to update img file
			fileRank := getRankFromFileName(items[i].ImagePath)
			if fileRank == items[i].Rank {
				fmt.Println("Item rank and path matches already, no swaps needed")
				continue
			}

			item := &items[i]

			fmt.Printf("Filename update req, filename is %d, but rank is %d\n", fileRank, item.Rank)
			swapItem := findSwapItem(items, item)

			if swapItem == nil {
				fmt.Println("no swap match found")
				continue
			}

			updateImage(item, swapItem, listID)

		}
	}

	for i := range items {
		fmt.Println(items[i].Title)
		fmt.Println(items[i].Rank)
		fmt.Println(items[i].ImagePath)
	}
	return items, nil
}

func getRankFromFileName(path string) int {
	splitPath := strings.Split(path, ".")
	if len(splitPath) < 2 {
		return -1
	}
	fileName := splitPath[0]
	fileRank, err := strconv.Atoi(fileName)
	if err != nil {
		return -1
	}
	return fileRank
}
